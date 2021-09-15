﻿// Copyright(c) 2020 Bitcoin Association.
// Distributed under the Open BSV software license, see the accompanying file LICENSE

using Dapper;
using MerchantAPI.APIGateway.Domain.Models;
using MerchantAPI.APIGateway.Domain.Repositories;
using MerchantAPI.Common.Clock;
using MerchantAPI.Common.Tasks;
using Microsoft.Extensions.Configuration;
using NBitcoin;
using Npgsql;
using System;
using System.Collections.Generic;
using System.Linq;
using MerchantAPI.Common.Json;

namespace MerchantAPI.APIGateway.Infrastructure.Repositories
{
  public class NodeRepositoryPostgres : INodeRepository
  {

    private readonly string connectionString;
    private static readonly Dictionary<string, Node> cache = new Dictionary<string, Node>();
    private readonly IClock clock;


    public NodeRepositoryPostgres(IConfiguration configuration, IClock clock)
    {
      connectionString = configuration["ConnectionStrings:DBConnectionString"];
      this.clock = clock ?? throw new ArgumentNullException(nameof(clock));
    }

    private void EnsureCache()
    {
      lock (cache)
      {
        if (!cache.Any())
        {
          foreach (var node in GetNodesDb())
          {
            cache.Add(GetCacheKey(node.ToExternalId()), node);
          }
        }
      }
    }

    private string GetCacheKey(string cachedKey)
    {
      return $"{cachedKey.ToLower()}";
    }


    public Node CreateNode(Node node)
    {
      EnsureCache();
      lock (cache)
      {
        var cachedKey = GetCacheKey(node.ToExternalId());
        if (cache.ContainsKey(cachedKey))
        {
          return null;
        }
        var createdNode = CreateNodeDb(node);
        if (createdNode != null)
        {
          cache.Add(cachedKey, createdNode);
        }
        return createdNode;
      }
    }

    private Node CreateNodeDb(Node node)
    {
      using var connection = new NpgsqlConnection(connectionString);
      RetryUtils.Exec(() => connection.Open());
      using var transaction = connection.BeginTransaction();

      string insertOrUpdate =
        "INSERT INTO Node " +
        "  (host, port, username, password, remarks, zmqnotificationsendpoint, nodestatus) " +
        "  VALUES (@host, @port, @username, @password, @remarks, @zmqnotificationsendpoint, @nodestatus)" +
        "  ON CONFLICT (host, port) DO NOTHING " +
        "  RETURNING nodeid as id, host, port, username, password, remarks, zmqnotificationsendpoint, nodestatus as status, lastError, lastErrorAt"
      ;

      var now = clock.UtcNow();

      var insertedNode = connection.Query<Node>(insertOrUpdate,
        new
        {
          host = node.Host.ToLower(),
          port = node.Port,
          username = node.Username,
          password = node.Password,
          remarks = node.Remarks,
          zmqnotificationsendpoint = node.ZMQNotificationsEndpoint,
          nodestatus = node.Status
        },
        transaction
      ).SingleOrDefault();
      transaction.Commit();

      return insertedNode;
    }

    public bool UpdateNode(Node node)
    {
      return UpdateNode(node, UpdateNodeDb);
    }

    private bool UpdateNode(Node node, Func<Node, (Node, bool)> func) 
    {
      EnsureCache();
      lock (cache)
      {
        var cachedKey = GetCacheKey(node.ToExternalId());
        if (!cache.ContainsKey(cachedKey))
        {
          return false;
        }
        (Node updatedNode, bool success) = func(node);
        if (success)
        {
          cache[cachedKey] = updatedNode;
        }
        return success;
      }
    }

    private (Node, bool) UpdateNodeDb(Node node)
    {
      using var connection = new NpgsqlConnection(connectionString);
      RetryUtils.Exec(() => connection.Open());
      using var transaction = connection.BeginTransaction();
      string update =
      "UPDATE Node " +
      "  SET  username=@username, password=@password, remarks=@remarks, zmqnotificationsendpoint=@zmqnotificationsendpoint " +
      "  WHERE host=@host AND port=@port" +
      "  RETURNING nodeid as id, host, port, username, password, remarks, zmqnotificationsendpoint, nodestatus as status, lastError, lastErrorAt";

      Node updatedNode = connection.Query<Node>(update,
        new
        {
          host = node.Host.ToLower(),
          port = node.Port,
          username = node.Username,
          password = node.Password,
          //nodestatus = node.Status, // NodeStatus is not present in ViewModel
          remarks = node.Remarks,
          zmqnotificationsendpoint = node.ZMQNotificationsEndpoint
        },
        transaction
      ).SingleOrDefault();
      transaction.Commit();

      return (updatedNode, updatedNode != null);
    }

    public bool UpdateNodeError(Node node)
    {
      return UpdateNode(node, UpdateNodeErrorDb);
    }

    private (Node, bool) UpdateNodeErrorDb(Node node)
    {
      using var connection = new NpgsqlConnection(connectionString);
      RetryUtils.Exec(() => connection.Open());
      using var transaction = connection.BeginTransaction();
      string update =
      "UPDATE Node " +
      "  SET  lastError=@lastError, lastErrorAt=@lastErrorAt " +
      "  WHERE nodeId=@nodeId" +
      "  RETURNING nodeid as id, host, port, username, password, remarks, zmqnotificationsendpoint, nodestatus as status, lastError, lastErrorAt";

      Node updatedNode = connection.Query<Node>(update,
        new
        {
          lastError = node.LastError,
          lastErrorAt = node.LastErrorAt,
          nodeId = node.Id
        },
        transaction
      ).SingleOrDefault();
      transaction.Commit();

      return (updatedNode, updatedNode != null);
    }


    public Node GetNode(string hostAndPort)
    {
      EnsureCache();
      lock (cache)
      {
        var cachedKey = GetCacheKey(hostAndPort);
        if (!cache.ContainsKey(cachedKey))
        {
          return null;
        }
        return cache.TryGet(cachedKey);
      }
    }


    public int DeleteNode(string hostAndPort)
    {
      EnsureCache();
      lock (cache)
      {
        var cachedKey = GetCacheKey(hostAndPort);
        if (!cache.ContainsKey(cachedKey))
        {
          return 0;
        }
        var deleted = DeleteNodeDb(cachedKey);
        if (deleted > 0)
        {
          cache.Remove(cachedKey, out var removedNode);
        }
        return deleted;
      }
    }


    private int DeleteNodeDb(string hostAndPort)
    {
      var (host, port) = Node.SplitHostAndPort(hostAndPort);

      using var connection = new NpgsqlConnection(connectionString);
      RetryUtils.Exec(() => connection.Open());
      using var transaction = connection.BeginTransaction();
      string cmd = "DELETE FROM Node WHERE host = @host AND  port = @port;";
      var result = connection.Execute(cmd,
        new
        {
          host = host.ToLower(),
          port
        },
        transaction
      );
      transaction.Commit();
      return result;
    }

    public IEnumerable<Node> GetNodes()
    {
      EnsureCache();

      lock (cache)
      {
        return cache.Values.ToArray();
      }
    }
    private IEnumerable<Node> GetNodesDb()
    {
      using var connection = new NpgsqlConnection(connectionString);
      RetryUtils.Exec(() => connection.Open());
      using var transaction = connection.BeginTransaction();
      string cmdText =
        @"SELECT nodeId as id, host, port, username, password, remarks, zmqnotificationsendpoint, nodeStatus as status, lastError, lastErrorAt FROM node ORDER by host, port";
      return connection.Query<Node>(cmdText, null, transaction);
    }

    public static void EmptyRepository(string connectionString)
    {
      using var connection = new NpgsqlConnection(connectionString);
      RetryUtils.Exec(() => connection.Open());
      string cmdText =
        "TRUNCATE node";
      connection.Execute(cmdText, null);

      lock (cache)
      {
        cache.Clear();
      }
    }
  }
}
