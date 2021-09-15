﻿// Copyright(c) 2020 Bitcoin Association.
// Distributed under the Open BSV software license, see the accompanying file LICENSE

using System;
using System.Threading;
using System.Threading.Tasks;
using MerchantAPI.APIGateway.Domain.Models;
using MerchantAPI.APIGateway.Domain.Models.Events;
using MerchantAPI.Common.Clock;
using MerchantAPI.Common.EventBus;
using Microsoft.Extensions.Logging;

namespace MerchantAPI.APIGateway.Domain.Actions
{
  public class BlockChainInfo : BackgroundServiceWithSubscriptions<BlockChainInfo>, IBlockChainInfo    
  {
    // Refresh every 60 seconds even if no ZMQ notification was received
    const int RefreshIntervalSeconds = 60;


    object objLock = new object();
    DateTime lastRefreshedAt;
    BlockChainInfoData cachedBlockChainInfo;
    IRpcMultiClient rpcMultiClient;
    private readonly IClock clock;

    EventBusSubscription<NewBlockDiscoveredEvent> newBlockDiscoveredSubscription;
    public BlockChainInfo(IRpcMultiClient rpcMultiClient, ILogger<BlockChainInfo> logger, IEventBus eventBus, IClock clock)
      : base(logger, eventBus)
    {
      this.rpcMultiClient= rpcMultiClient?? throw new ArgumentNullException(nameof(rpcMultiClient));
      this.clock = clock ?? throw new ArgumentNullException(nameof(clock));
      lastRefreshedAt = clock.UtcNow();
    }
    public BlockChainInfoData GetInfo()
    {

      lock (objLock)
      {
        // Refresh if needed
        if (cachedBlockChainInfo == null || (clock.UtcNow() - lastRefreshedAt).TotalSeconds > RefreshIntervalSeconds)
        {
          var blockChainInfoTask = rpcMultiClient.GetWorstBlockchainInfoAsync();
          var networkInfoTask = rpcMultiClient.GetAnyNetworkInfoAsync();

          // Note that the following call will block.
          Task.WhenAll(blockChainInfoTask, networkInfoTask).GetAwaiter().GetResult();

          cachedBlockChainInfo = new BlockChainInfoData(
            blockChainInfoTask.Result.BestBlockHash,
            blockChainInfoTask.Result.Blocks,
            new ConsolidationTxParameters(networkInfoTask.Result)
          );
          lastRefreshedAt = clock.UtcNow();
        }

        return cachedBlockChainInfo;
      }
    }

    protected override void UnsubscribeFromEventBus()
    {
      eventBus?.TryUnsubscribe(newBlockDiscoveredSubscription);
      newBlockDiscoveredSubscription = null;
    }


    protected override void SubscribeToEventBus(CancellationToken stoppingToken)
    {
      newBlockDiscoveredSubscription = eventBus.Subscribe<NewBlockDiscoveredEvent>();

      _ = newBlockDiscoveredSubscription.ProcessEventsAsync(stoppingToken, logger, NewBlockDiscoveredAsync);
    }

    private Task NewBlockDiscoveredAsync(NewBlockDiscoveredEvent arg)
    {
      lock (objLock)
      {
        lastRefreshedAt = DateTime.MinValue;
        // Note that RpcMultiClient.GetBlockchainInfoAsync will return the WORST block from all nodes
        // so in the case of X nodes reporting the best block, we will do actually do X^2 GetBlockchainInfoAsync
        // calls and only when the last node will catchup, GetBlockchainInfoAsync will report changes result.
        // We could optimize this by tracking  per-node state in this class. This would also require
        // that we subscribe to Node integration events.
      }

      return Task.CompletedTask;
    }

    protected override Task ProcessMissedEvents()
    {
      // Nothing to do here
      return Task.CompletedTask;

    }
  }
}
