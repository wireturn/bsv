﻿// Copyright(c) 2020 Bitcoin Association.
// Distributed under the Open BSV software license, see the accompanying file LICENSE

namespace MerchantAPI.APIGateway.Domain.Models
{
  public class SubmitTransactionOneResponse
  {
    public string Txid { get; set; }

    public string ReturnResult { get; set; }

    public string ResultDescription { get; set; }

    public SubmitTransactionConflictedTxResponse[] ConflictedWith { get; set; }

  }
}
