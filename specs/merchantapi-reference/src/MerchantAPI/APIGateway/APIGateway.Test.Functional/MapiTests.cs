﻿// Copyright(c) 2020 Bitcoin Association.
// Distributed under the Open BSV software license, see the accompanying file LICENSE

using MerchantAPI.APIGateway.Domain;
using MerchantAPI.APIGateway.Domain.Models;
using MerchantAPI.APIGateway.Domain.ViewModels;
using MerchantAPI.APIGateway.Rest.ViewModels;
using MerchantAPI.APIGateway.Test.Functional.Mock;
using MerchantAPI.APIGateway.Test.Functional.Server;
using MerchantAPI.APIGateway.Test.Functional.Attributes;
using MerchantAPI.Common.Test.Clock;
using MerchantAPI.Common.Json;
using Microsoft.AspNetCore.TestHost;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using NBitcoin;
using NBitcoin.Altcoins;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Net;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Net.Mime;
using System.Threading.Tasks;

namespace MerchantAPI.APIGateway.Test.Functional
{

  [TestClass]
  public class MapiTests : TestBase
  {
    public TestContext TestContext { get; set; }

    void AddMockNode(int nodeNumber)
    {
      var mockNode = new Node(0, "mockNode" + nodeNumber, 0, "mockuserName", "mockPassword", "This is a mock node #" + nodeNumber,
        null, (int)NodeStatus.Connected, null, null);

      _ = Nodes.CreateNodeAsync(mockNode).Result;
    }

    [TestInitialize]
    public void TestInitialize()
    {
      //Retrive OverrideSettingAttribute data (setting name and value)
      List<KeyValuePair<string, string>> overridenSettings = new List<KeyValuePair<string, string>>();     
      var overrideSettingsAttributes = GetType().GetMethod(TestContext.TestName).GetCustomAttributes(true).Where(a => a.GetType() == typeof(OverrideSettingAttribute));
      foreach (var attribute in overrideSettingsAttributes)
      {
        OverrideSettingAttribute overrideSettingsAttribute = (OverrideSettingAttribute)attribute;
        overridenSettings.Add(new KeyValuePair<string, string>(overrideSettingsAttribute.SettingName, overrideSettingsAttribute.SettingValue.ToString()));
      }

      Initialize(mockedServices: true, overridenSettings);
      AddMockNode(0);
    }

    public override TestServer CreateServer(bool mockedServices, TestServer serverCallback, string dbConnectionString, IEnumerable<KeyValuePair<string, string>> overridenSettings = null)
    {
      return new TestServerBase(DbConnectionStringDDL).CreateServer<MapiServer, APIGatewayTestsMockStartup, APIGatewayTestsStartup>(mockedServices, serverCallback, dbConnectionString, overridenSettings);
    }


    [TestCleanup]
    public void TestCleanup()
    {
      Cleanup();
    }

    static void VerifySignature((SignedPayloadViewModel response, HttpResponseMessage httpResponse) response)
    {
      Assert.IsTrue(JsonEnvelopeSignature.VerifySignature(response.response.ToDomainObject()), "Signature is invalid");
    }

    void AssertIsOK(SubmitTransactionResponseViewModel response, string expectedTxId, string expectedResult = "success", string expectedDescription = "")
    {

      Assert.AreEqual("1.3.0", response.ApiVersion);
      Assert.IsTrue((MockedClock.UtcNow - response.Timestamp).TotalSeconds < 60);
      Assert.AreEqual(expectedResult, response.ReturnResult);
      Assert.AreEqual(expectedDescription, response.ResultDescription);

      Assert.AreEqual(MinerId.GetCurrentMinerIdAsync().Result, response.MinerId);
      Assert.AreEqual(BlockChainInfo.GetInfo().BestBlockHeight, response.CurrentHighestBlockHeight);
      Assert.AreEqual(BlockChainInfo.GetInfo().BestBlockHash, response.CurrentHighestBlockHash);
      Assert.AreEqual(expectedTxId, response.Txid);
    }
    void AssertIsOK(FeeQuoteViewModelGet response)
    {

      Assert.AreEqual("1.3.0", response.ApiVersion);
      Assert.IsTrue((MockedClock.UtcNow - response.Timestamp).TotalSeconds < 60);

      Assert.AreEqual(MinerId.GetCurrentMinerIdAsync().Result, response.MinerId);
      Assert.AreEqual(BlockChainInfo.GetInfo().BestBlockHeight, response.CurrentHighestBlockHeight);
      Assert.AreEqual(BlockChainInfo.GetInfo().BestBlockHash, response.CurrentHighestBlockHash);
    }

    [TestMethod]
    public async Task GetFeeQuote()
    {
      var response = await Get<SignedPayloadViewModel>(
        client, MapiServer.ApiMapiQueryFeeQuote, HttpStatusCode.OK);

      var payload = response.ExtractPayload<FeeQuoteViewModelGet>();
      AssertIsOK(payload);

      using (MockedClock.NowIs(DateTime.UtcNow.AddMinutes(FeeQuoteRepositoryMock.quoteExpiryMinutes + 1)))
      {
        // should return same
        response = await Get<SignedPayloadViewModel>(
          client, MapiServer.ApiMapiQueryFeeQuote, HttpStatusCode.OK);

        payload = response.ExtractPayload<FeeQuoteViewModelGet>();
        AssertIsOK(payload);
      }

    }


    [TestMethod]
    public async Task GetFeeQuoteAuthenticated()
    {
      RestAuthentication = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI1IiwibmJmIjoxNTk5NDExNDQzLCJleHAiOjE5MTQ3NzE0NDMsImlhdCI6MTU5OTQxMTQ0MywiaXNzIjoiaHR0cDovL215c2l0ZS5jb20iLCJhdWQiOiJodHRwOi8vbXlhdWRpZW5jZS5jb20ifQ.Z43NASAbIxMZrL2MzbJTJD30hYCxhoAs-8heDjQMnjM";
      (SignedPayloadViewModel response, HttpResponseMessage httpResponse) response = await GetWithHttpResponseReturned<SignedPayloadViewModel>(
                     client, MapiServer.ApiMapiQueryFeeQuote, HttpStatusCode.NotFound);
      Assert.AreEqual("Not Found", response.httpResponse.ReasonPhrase);

      feeQuoteRepositoryMock.FeeFileName = "feeQuotesWithIdentity.json";
      response = await GetWithHttpResponseReturned<SignedPayloadViewModel>(
                 client, MapiServer.ApiMapiQueryFeeQuote, HttpStatusCode.OK);
      var payload = response.response.ExtractPayload<FeeQuoteViewModelGet>();
      AssertIsOK(payload);
    }

    [TestMethod]
    public async Task GetFeeQuote_WithInvalidAuthentication()
    {
      feeQuoteRepositoryMock.FeeFileName = "feeQuotesWithIdentity.json";
      // TokenManager.exe generate -n testName -i http://mysite.com -a http://myaudience.com -k thisisadevelopmentkey -d 3650
      var ValidRestAuthentication = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI1IiwibmJmIjoxNTk5NDExNDQzLCJleHAiOjE5MTQ3NzE0NDMsImlhdCI6MTU5OTQxMTQ0MywiaXNzIjoiaHR0cDovL215c2l0ZS5jb20iLCJhdWQiOiJodHRwOi8vbXlhdWRpZW5jZS5jb20ifQ.Z43NASAbIxMZrL2MzbJTJD30hYCxhoAs-8heDjQMnjM";
     
      RestAuthentication = ValidRestAuthentication+"invalid";
      var response = await Get<SignedPayloadViewModel>(
                 client, MapiServer.ApiMapiQueryFeeQuote, HttpStatusCode.Unauthorized);
      Assert.IsNull(response);
    }

    [TestMethod]
    public async Task GetFeeQuote_TestBearerToken()
    {
      feeQuoteRepositoryMock.FeeFileName = "feeQuotesWithIdentity.json";
      // test authentication: same provider and identity as defined in json - should succeed
      // TokenManager.exe generate -n testName -i http://mysite.com -a http://myaudience.com -k thisisadevelopmentkey -d 3650
      RestAuthentication = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI1IiwibmJmIjoxNTk5NDExNDQzLCJleHAiOjE5MTQ3NzE0NDMsImlhdCI6MTU5OTQxMTQ0MywiaXNzIjoiaHR0cDovL215c2l0ZS5jb20iLCJhdWQiOiJodHRwOi8vbXlhdWRpZW5jZS5jb20ifQ.Z43NASAbIxMZrL2MzbJTJD30hYCxhoAs-8heDjQMnjM";
      var response = await Get<SignedPayloadViewModel>(
                 client, MapiServer.ApiMapiQueryFeeQuote, HttpStatusCode.OK);
      var payload = response.ExtractPayload<FeeQuoteViewModelGet>();
      AssertIsOK(payload);

      // different user, same provider, same authority - should succeed
      // TokenManager.exe generate -n testName -i http://mysite.com -a http://myaudience.com -k thisisadevelopmentkey -d 3650
      RestAuthentication = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0TmFtZSIsIm5iZiI6MTYwMzg2NjAyOCwiZXhwIjoxOTE5MjI2MDI4LCJpYXQiOjE2MDM4NjYwMjgsImlzcyI6Imh0dHA6Ly9teXNpdGUuY29tIiwiYXVkIjoiaHR0cDovL215YXVkaWVuY2UuY29tIn0.01Rm6t4GBScDwgoOnFwBjjvgu6U5YBK7qlCTg-_BF6c";
      response = await Get<SignedPayloadViewModel>(
                 client, MapiServer.ApiMapiQueryFeeQuote, HttpStatusCode.OK);
      payload = response.ExtractPayload<FeeQuoteViewModelGet>();
      AssertIsOK(payload);

      // same user, different (invalid) provider, same authority - should fail
      //TokenManager.exe generate -n testName - i http://test.com -a http://myaudience.com -k thisisadevelopmentkey -d 3650
      RestAuthentication = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0TmFtZSIsIm5iZiI6MTYwMzg2NjQ4OCwiZXhwIjoxOTE5MjI2NDg4LCJpYXQiOjE2MDM4NjY0ODgsImlzcyI6Imh0dHA6Ly90ZXN0LmNvbSIsImF1ZCI6Imh0dHA6Ly9teWF1ZGllbmNlLmNvbSJ9.oGxXXbTj0yUf0UrwOF44bbRMt-Xe6YjAyuy4A3jrbbU";
      response = await Get<SignedPayloadViewModel>(
           client, MapiServer.ApiMapiQueryFeeQuote, HttpStatusCode.Unauthorized);
      Assert.IsNull(response);

      // same user and provider, different authority
      // TokenManager.exe generate -n 5 -i http://mysite.com -a http://testaudience.com -k thisisadevelopmentkey -d 3650
      RestAuthentication = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI1IiwibmJmIjoxNjAzODY2NzAxLCJleHAiOjE5MTkyMjY3MDEsImlhdCI6MTYwMzg2NjcwMSwiaXNzIjoiaHR0cDovL215c2l0ZS5jb20iLCJhdWQiOiJodHRwOi8vdGVzdGF1ZGllbmNlLmNvbSJ9.d0TU7em4_8ZzO8A3YGxVwyl0ElpDQIu35auPSa24i48";
      response = await Get<SignedPayloadViewModel>(
     client, MapiServer.ApiMapiQueryFeeQuote, HttpStatusCode.Unauthorized);
      Assert.IsNull(response);
    }

    [TestMethod]
    public async Task SubmitTransaction_WithInvalidAuthentication()
    {
      feeQuoteRepositoryMock.FeeFileName = "feeQuotesWithIdentity.json";
      var ValidRestAuthentication = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI1IiwibmJmIjoxNTk5NDExNDQzLCJleHAiOjE5MTQ3NzE0NDMsImlhdCI6MTU5OTQxMTQ0MywiaXNzIjoiaHR0cDovL215c2l0ZS5jb20iLCJhdWQiOiJodHRwOi8vbXlhdWRpZW5jZS5jb20ifQ.Z43NASAbIxMZrL2MzbJTJD30hYCxhoAs-8heDjQMnjM";
      RestAuthentication = ValidRestAuthentication + "invalid";

      var txBytes = HelperTools.HexStringToByteArray(txC3Hex);

      var reqContent = new ByteArrayContent(txBytes);
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Octet);

      var response =
        await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.Unauthorized);
      Assert.IsNull(response.response);

      Assert.AreEqual(0, rpcClientFactoryMock.AllCalls.FilterCalls("mocknode0:sendrawtransactions/").Count()); // no calls

    }

    [TestMethod]
    public async Task SubmitTransactionBinary()
    {
      var txBytes = HelperTools.HexStringToByteArray(txC3Hex);

      var reqContent = new ByteArrayContent(txBytes);
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Octet);

      var response =
        await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);

      rpcClientFactoryMock.AllCalls.AssertContains("mocknode0:sendrawtransactions/", "mocknode0:sendrawtransactions/" + txC3Hash);

      var payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();

      // Check if all fields are set
      AssertIsOK(payload, txC3Hash);
    }

    [TestMethod]
    public async Task SubmitTransactionDuplicateError()
    {
      var txBytes = HelperTools.HexStringToByteArray(txC3Hex);

      var reqContent = new ByteArrayContent(txBytes);
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Octet);

      var response =
        await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);

      rpcClientFactoryMock.AllCalls.AssertContains("mocknode0:sendrawtransactions/", "mocknode0:sendrawtransactions/" + txC3Hash);

      var payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();

      // Check if all fields are set
      AssertIsOK(payload, txC3Hash);

      response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);
      payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();
      Assert.AreEqual(payload.ReturnResult, "failure");
      Assert.AreEqual(payload.ResultDescription, "Transaction already known");
    }

    [TestMethod]
    public async Task SubmitTransactionJson()
    {
      var reqContent = new StringContent($"{{ \"rawtx\": \"{txC3Hex}\" }}");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);

      var response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);

      rpcClientFactoryMock.AllCalls.AssertContains("mocknode0:sendrawtransactions/", "mocknode0:sendrawtransactions/" + txC3Hash);

      var payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();

      // Check if all fields are set
      AssertIsOK(payload, txC3Hash);
      Assert.AreEqual("", payload.ResultDescription); // Description should be "" (not null)
    }

    private int GetBytesForScriptLength(ulong totalBytes)
    {
      if (totalBytes < byte.MaxValue) // uint8 == byte
      {
        return 1;
      }
      else if (totalBytes < UInt16.MaxValue) // if script length can not be encoded in single byte we need additional data
      {
        return 3; // saved as variable length integer (0xFD followed by the length as uint16_t)
      }
      else if (totalBytes < UInt32.MaxValue)
      {
        return 5;
      }
      else if (totalBytes < UInt64.MaxValue)
      {
        return 9;
      }
      else
      {
        throw new ArgumentException("Script is too big.");
      }
    }


    /// <summary>
    /// Create a new transaction with is totalBytes long. Out of this totalBytes, dataBytes are spend for 
    /// </summary>
    /// <param name="fundingTx">Input transaction. It's first output will be used as funding for new transaction</param>
    /// <param name="totalBytes">Total desired length of created transaction</param>
    /// <param name="dataBytes">Number of data bytes (OP_FALSE transaction ....) that this transaction  should contain</param>
    /// <param name="totalFees"> total fees payed by this transaction</param>
    /// <returns></returns>
    Transaction CreateTransaction(Transaction fundingTx, long totalBytes, long dataBytes, long totalFees)
    {
      if (dataBytes > 0)
      {
        if (dataBytes < 2)
        {
          throw new ArgumentException($"nameof(dataBytes) should be at least 2, since script must start with OP_FALSE OP_RETURN");
        }
      }
      
      var remainingMoney = fundingTx.Outputs[0].Value - totalFees;
      if (remainingMoney < 0L)
      {
        throw new ArgumentException("Fee is too large (or funding output is to low)");
      }

      var tx = BCash.Instance.Regtest.CreateTransaction();
      tx.Inputs.Add(new TxIn(new OutPoint(fundingTx, 0)));

      long sizeOfSingleOutputWithoutScript = sizeof(ulong) + GetBytesForScriptLength((ulong) (totalBytes - dataBytes)); // 9+:	A list of 1 or more transaction outputs or destinations for coins
      long overHead =
           tx.ToBytes().Length // length of single input
          + dataBytes == 0 ? 0 : tx.ToBytes().Length + sizeOfSingleOutputWithoutScript;

      long normalBytes = totalBytes - dataBytes - overHead;

      if (normalBytes > 0 && dataBytes > 0) // Overhead also depends on number of outputs - if this is true we have two outputs 
      {
        normalBytes -= (sizeof(ulong) + GetBytesForScriptLength((ulong)dataBytes));
      }

      if (normalBytes > 0)
      {
        var scriptBytes = new byte[normalBytes];
        tx.Outputs.Add(new TxOut(remainingMoney, new Script(scriptBytes)));
        remainingMoney = 0L;
      }
      else if (normalBytes < 0)
      {
        throw new ArgumentException("Argument Databytes is too low.");
      }
      if (dataBytes > 0)
      {
        var scriptBytes = new byte[dataBytes];
        scriptBytes[0] = (byte)OpcodeType.OP_FALSE;
        scriptBytes[1] = (byte)OpcodeType.OP_RETURN;
        tx.Outputs.Add(new TxOut(remainingMoney, new Script(scriptBytes)));
      }

      Assert.AreEqual(totalBytes, tx.ToBytes().Length, "Failed to create transaction of desired length");

      return tx;
      
    }

    [TestMethod]
    public async Task SubmitTransactionJsonMinFee()
    {
      // register tx0
      rpcClientFactoryMock.SetUpTransaction(txC0Hex);
      var tx0 = HelperTools.ParseBytesToTransaction(HelperTools.HexStringToByteArray(txC0Hex));

      int txLength = 160;
      var fee = feeQuoteRepositoryMock.GetValidFeeQuotesByIdentity(null).Single().Fees.Single(x => x.FeeType == Const.FeeType.Standard);
      var minRequiredFees = Math.Min((txLength * fee.RelayFee.Satoshis) / fee.RelayFee.Bytes, // 40
                          (txLength * fee.MiningFee.Satoshis) / fee.MiningFee.Bytes); // 80

      var tx1 = CreateTransaction(tx0, txLength, 0, minRequiredFees); // submit tx1 should succeed

      var reqContent = new StringContent($"{{ \"rawtx\": \"{ tx1.ToHex() }\" }}");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);

      var response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);

      rpcClientFactoryMock.AllCalls.AssertContains("mocknode0:sendrawtransactions/", "mocknode0:sendrawtransactions/" + tx1.GetHash().ToString());

      var payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();

      // Check if all fields are set
      AssertIsOK(payload, tx1.GetHash().ToString());

    }

    [TestMethod]
    public async Task SubmitTransactionJsonMinSumFee()
    {
      // register tx0
      rpcClientFactoryMock.SetUpTransaction(txC0Hex);
      var tx0 = HelperTools.ParseBytesToTransaction(HelperTools.HexStringToByteArray(txC0Hex));
 
      int txLength = 160;
      var fee = feeQuoteRepositoryMock.GetValidFeeQuotesByIdentity(null).Single().Fees.Single(x => x.FeeType == Const.FeeType.Standard);
      var minRequiredFees = Math.Min((txLength * fee.RelayFee.Satoshis) / fee.RelayFee.Bytes, // 40
                          (txLength * fee.MiningFee.Satoshis) / fee.MiningFee.Bytes); // 80
      var tx1 = CreateTransaction(tx0, txLength, 0, minRequiredFees - 1); // submit tx1 should fail

      var reqContent = new StringContent($"{{ \"rawtx\": \"{ tx1.ToHex() }\" }}");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);

      var response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);

      Assert.AreEqual(0, rpcClientFactoryMock.AllCalls.FilterCalls("mocknode0:sendrawtransactions/").Count()); // no calls, to submit txs since we do not pay enough fee

      var payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();
      // Check if all fields are set
      AssertIsOK(payload, tx1.GetHash().ToString(), "failure", "Not enough fees");

      var tx2 = CreateTransaction(tx0, txLength, 0, minRequiredFees + 1); // submit tx2 should succeed

      reqContent = new StringContent($"{{ \"rawtx\": \"{ tx2.ToHex() }\" }}");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);

      response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);

      rpcClientFactoryMock.AllCalls.AssertContains("mocknode0:sendrawtransactions/", "mocknode0:sendrawtransactions/" + tx2.GetHash());

      payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();

      // Check if all fields are set
      AssertIsOK(payload, tx2.GetHash().ToString());
    }


    [TestMethod]
    public async Task SubmitTransactionJsonMinFeeTypeData()
    {
      // register tx0
      rpcClientFactoryMock.SetUpTransaction(txC0Hex);
      var tx0 = HelperTools.ParseBytesToTransaction(HelperTools.HexStringToByteArray(txC0Hex));

      long txLength = 160;
      long dataLength = 100;
      long standard = txLength-dataLength;
      var fee = feeQuoteRepositoryMock.GetValidFeeQuotesByIdentity(null).Single().Fees.Single(x => x.FeeType == Const.FeeType.Data);
      var minRequiredFees = Math.Min((dataLength * fee.RelayFee.Satoshis) / fee.RelayFee.Bytes, // 20
                    (dataLength * fee.MiningFee.Satoshis) / fee.MiningFee.Bytes); // 40
      fee = feeQuoteRepositoryMock.GetValidFeeQuotesByIdentity(null).Single().Fees.Single(x => x.FeeType == Const.FeeType.Standard);
      minRequiredFees += Math.Min((standard * fee.RelayFee.Satoshis) / fee.RelayFee.Bytes, // 15
              (standard * fee.MiningFee.Satoshis) / fee.MiningFee.Bytes); // 30


      var tx1 = CreateTransaction(tx0, txLength, dataLength, minRequiredFees); // submit tx1 should succeed

      var reqContent = new StringContent($"{{ \"rawtx\": \"{ tx1.ToHex() }\" }}");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);

      var response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);

      rpcClientFactoryMock.AllCalls.AssertContains("mocknode0:sendrawtransactions/", "mocknode0:sendrawtransactions/" + tx1.GetHash().ToString());

      var payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();

      // Check if all fields are set
      AssertIsOK(payload, tx1.GetHash().ToString());

    }

    public async Task SubmitTransactionJsonMinSumFeeTypeData()
    {
      // register tx0
      rpcClientFactoryMock.SetUpTransaction(txC0Hex);
      var tx0 = HelperTools.ParseBytesToTransaction(HelperTools.HexStringToByteArray(txC0Hex));

      long txLength = 160;
      long dataLength = 100;
      long standard = txLength - dataLength;
      var fee = feeQuoteRepositoryMock.GetValidFeeQuotesByIdentity(null).Single().Fees.Single(x => x.FeeType == Const.FeeType.Data);
      var minRequiredFees = Math.Min((dataLength * fee.RelayFee.Satoshis) / fee.RelayFee.Bytes, // 20
                    (dataLength * fee.MiningFee.Satoshis) / fee.MiningFee.Bytes); // 40
      fee = feeQuoteRepositoryMock.GetValidFeeQuotesByIdentity(null).Single().Fees.Single(x => x.FeeType == Const.FeeType.Standard);
      minRequiredFees += Math.Min((standard * fee.RelayFee.Satoshis) / fee.RelayFee.Bytes, // 15
              (standard * fee.MiningFee.Satoshis) / fee.MiningFee.Bytes); // 30

      var tx1 = CreateTransaction(tx0, txLength, 0, minRequiredFees - 1); // submit tx1 should fail

      var reqContent = new StringContent($"{{ \"rawtx\": \"{ tx1.ToHex() }\" }}");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);

      var response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);

      Assert.AreEqual(0, rpcClientFactoryMock.AllCalls.FilterCalls("mocknode0:sendrawtransactions/").Count()); // no calls, to submit txs since we do not pay enough fee

      var payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();
      // Check if all fields are set
      AssertIsOK(payload, tx1.GetHash().ToString(), "failure", "Not enough fees");


      var tx2 = CreateTransaction(tx0, txLength, 0, minRequiredFees + 1); // submit tx2 should succeed

      reqContent = new StringContent($"{{ \"rawtx\": \"{ tx2.ToHex() }\" }}");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);

      response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);

      rpcClientFactoryMock.AllCalls.AssertContains("mocknode0:sendrawtransactions/", "mocknode0:sendrawtransactions/" + tx2.GetHash().ToString());

      payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();

      // Check if all fields are set
      AssertIsOK(payload, tx2.GetHash().ToString());

    }

    [TestMethod]
    public async Task SubmitTransactionJsonMinFeeTypeDataTestBigLength()
    {
      // register tx0
      rpcClientFactoryMock.SetUpTransaction(txC0Hex);
      var tx0 = HelperTools.ParseBytesToTransaction(HelperTools.HexStringToByteArray(txC0Hex));

      long txLength = 400000;
      long dataLength = 350000;
      long standard = txLength - dataLength;
      var fee = feeQuoteRepositoryMock.GetValidFeeQuotesByIdentity(null).Single().Fees.Single(x => x.FeeType == Const.FeeType.Data);
      var minRequiredFees = Math.Min((dataLength * fee.RelayFee.Satoshis) / fee.RelayFee.Bytes, 
                    (dataLength * fee.MiningFee.Satoshis) / fee.MiningFee.Bytes); 
      fee = feeQuoteRepositoryMock.GetValidFeeQuotesByIdentity(null).Single().Fees.Single(x => x.FeeType == Const.FeeType.Standard);
      minRequiredFees += Math.Min((standard * fee.RelayFee.Satoshis) / fee.RelayFee.Bytes, 
              (standard * fee.MiningFee.Satoshis) / fee.MiningFee.Bytes); 


      var tx1 = CreateTransaction(tx0, txLength, dataLength, minRequiredFees); // submit tx1 should succeed

      var reqContent = new StringContent($"{{ \"rawtx\": \"{ tx1.ToHex() }\" }}");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);

      var response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);

      rpcClientFactoryMock.AllCalls.AssertContains("mocknode0:sendrawtransactions/", "mocknode0:sendrawtransactions/" + tx1.GetHash().ToString());

      var payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();

      // Check if all fields are set
      AssertIsOK(payload, tx1.GetHash().ToString());

    }


    [TestMethod]
    [DataRow("something", HttpStatusCode.OK, "CallbackUrl: something should be a valid URL")]
    [DataRow("invalidScheme://www.something.com", HttpStatusCode.OK, "CallbackUrl: invalidScheme://www.something.com uses invalid scheme. Only 'http' and 'https' are supported")]
    [DataRow("http://www.something.com", HttpStatusCode.OK, "")]
    [DataRow("https://www.something.com", HttpStatusCode.OK, "")]
    public async Task SubmitTransactionJsonInvalidCallbackUrl(string url, HttpStatusCode expectedCode, string returnResult)
    {
      // Taken from whatsonchain
      // tx0Hex -> tx1Hex -> tx2Hex -> tx3Hex -> tx4Hex -> tx5Hex -> tx6Hex
      var tx0Hex = "0100000001af9c50faf7863139e5e2f2350a0f8fe130271562008369ccc53b63cc7eb2e1ff000000006b483045022100def69a937e05042cf3c72819e1dc6be8aa450bde1d91f8280ebafdcf369d6f8302205c303b9e7b5f18678e498266fc4405f71df174fd1632bccff11bb19283777cb44121028b85a286cbc1eefc0dd68b350d09d9942ec58dbbfc1c320de3d762491e9ff84bffffffff0240420f00000000001976a9146483af0f192d1e504df638e1775084b0b01dc14c88acc4aba254010000001976a9146db8ba46c120dd90e62f85ddda29207cfe8309b888ac00000000";
      var tx1Hex = "01000000010eb4f43b8909b8c6c0172cff022f5923adb6c87f5f52f942667aefc39e2fb6a7010000006b483045022100c5a29893ebea43f1afe14988b3a89aa81b3eb811ddc0adb249de52b97dcf330f0220411cddc31309c8461f45a4b33f5a184d5098880d601f6bc3942bb51f64cbd3674121020eea82bfd273a2588a0e737e6755a8169504c00dd46326cef5a08810761557a3ffffffff02a2777153010000001976a91497847593587fa8e15f0177a8fba47a356230a3cf88ac002d3101000000001976a914d33ed200d58a7ff330b1ca8af2a022ae42555c8588ac00000000";
      var tx2Hex = "01000000014cb8f7c1091fba0a09ccb1073c5c41ec5e1bba114502199b82481d2c0d8ef8b1000000006b483045022100ac5912a8f22a822216f7bfa177c0f5a8fefa287ededd7d108aa7d73bee6ed738022023bef4bfc0412d3143f4579f4e91536dc0044407dd0d028ec3bedb0dda2867c741210283816272e909371d5c006d4aaab6b0fa26ec76c85190abaf7a88db23b7b34a89ffffffff02800af270000000001976a914d2c2f9ebcabbdca1ff254ff0c21dc3937a863f1d88ac00667fe2000000001976a91488917aa35b9e01c9e40f00897d4d54de0d825a4288ac00000000";
      var tx3Hex = "01000000015760c1f6d67b58d13ffc1d72ab8f518f86499cf21c5284097f2f1a0e71335b7a000000006a47304402202f496eede845ef3b484b23ea256662c991d5ff395c209c65163d32acc86d134902207bad37320f3591e2f224accfd2c0dad1d57dc6bd110a748bd90494b7f32206c24121032de3c6f4e19d9970548abbd395947a9e3466a6b6f98f4b05ba8fa744975d7bebffffffff02e353a733000000001976a914fcdee390ba63b52808fc5f6fea227eb2cf1114e188ac7baf4a3d000000001976a9147d397b7af65b6aa506a6359154784d46676e804888ac00000000";
      var tx4Hex = "0100000001495d68ba120b956694e9450b25d1d022ca496507f5b0539640acbee86c8e1422010000006b483045022100eb50b584b5d0a232241eac80a288a8c5e8b1599291063aaea4c6a358608476bf022063ed106993cce098d069a750570a860cea297b24ba712f550712dc1efbd64b84412102eacb5afbc0d9609d7d92815960ac637695bb83b052c89703659cda7c2739e4faffffffff0200c2eb0b000000001976a914ddff71008839ec12ab8d66f74c0454f36b01438e88ac38f55d31000000001976a914c713b6d656411e584cc5fc3125a43073c3175e2488ac00000000";
      var tx5Hex = "01000000013aea7ffea18d6558f36913c5bf80dc3386ee09786be9f6cbe37c0a09414f090d010000006b483045022100d47a68d211ab71f4a128550e383598efa92551a93dd8df666b74acd545155d9b02201e0f842a8f349673d21d61e0a550427a3109be7cbd2bf9ccc27ae562e71c681141210295165fe23f09025fa9ca6e59da2025949045de4f3370110773fc4029f8185075ffffffff020084d717000000001976a9144afe74d11e424d38d13f822167476d10a83379dd88accd9b8519000000001976a914a99b0b0bd8ce2846f56c8cfd566141dfa55ffacd88ac00000000";
      var tx6Hex = "0100000001ecaa061ebcd099bdb3ad8417a7d18467a4ad87208aaf31552505f310780f3046010000006a47304402200b5079f89b4293e3a6ea9d0e269b23751872afd5741c7ad314c986556e5eaecb0220273494b6e8532126c8c4c33b566be6dff5c5257dfc600840000be81d0d73f0b04121033b87fb2d1d6285439c7b3164ec9f1f8fcf0929c29c44bce7cc58ac46b14242b500000000014d9a8519000000001976a914dbe02557c534a026868b253da20c8eb3956dbb8788ac00000000";

      rpcClientFactoryMock.AddKnownTransaction(HelperTools.HexStringToByteArray(tx0Hex));
      rpcClientFactoryMock.AddKnownTransaction(HelperTools.HexStringToByteArray(tx1Hex));
      rpcClientFactoryMock.AddKnownTransaction(HelperTools.HexStringToByteArray(tx2Hex));
      rpcClientFactoryMock.AddKnownTransaction(HelperTools.HexStringToByteArray(tx3Hex));
      rpcClientFactoryMock.AddKnownTransaction(HelperTools.HexStringToByteArray(tx4Hex));
      rpcClientFactoryMock.AddKnownTransaction(HelperTools.HexStringToByteArray(tx5Hex));
      rpcClientFactoryMock.AddKnownTransaction(HelperTools.HexStringToByteArray(tx6Hex));
      // Test submitting single tx through JSON
      var reqContent = new StringContent($"{{ \"rawtx\": \"{tx1Hex}\", \"callbackUrl\" : \"{url}\" }}");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);
      var resp = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, expectedCode);
      var txRespViewModel = HelperTools.JSONDeserialize<SubmitTransactionResponseViewModel>(resp.response.Payload);
      Assert.AreEqual(returnResult, txRespViewModel.ResultDescription);

      // Test submitting multiple txs through JSON, include callback in URL
      reqContent = new StringContent($"[ {{ \"rawtx\": \"{tx2Hex}\", \"callbackUrl\" : \"{url}\" }}]");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);
      resp = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransactions, client, reqContent, expectedCode);
      var txsRespViewModel = HelperTools.JSONDeserialize<SubmitTransactionsResponseViewModel>(resp.response.Payload);
      Assert.AreEqual(returnResult, txsRespViewModel.Txs[0].ResultDescription);

      // Test submitting multiple txs through JSON, use default callback
      reqContent = new StringContent($"[ {{ \"rawtx\": \"{tx3Hex}\" }}]");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);
      resp = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransactions+ $"?defaultCallbackUrl={url}", client, reqContent, expectedCode);
      txsRespViewModel = HelperTools.JSONDeserialize<SubmitTransactionsResponseViewModel>(resp.response.Payload);
      Assert.AreEqual(returnResult, txsRespViewModel.Txs[0].ResultDescription);


      // Test submitting multiple txs through JSON
      reqContent = new StringContent($"[ {{ \"rawtx\": \"{tx4Hex}\", \"callbackUrl\" : \"{url}\" }}]");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);
      resp = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransactions, client, reqContent, expectedCode);
      txsRespViewModel = HelperTools.JSONDeserialize<SubmitTransactionsResponseViewModel>(resp.response.Payload);
      Assert.AreEqual(returnResult, txsRespViewModel.Txs[0].ResultDescription);

      // Test submitting single transaction through Binary
      var txBytes = HelperTools.HexStringToByteArray(tx5Hex);

      var reqContentBin = new ByteArrayContent(txBytes);
      reqContentBin.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Octet);

      resp = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction+$"?callbackUrl={url}", client, reqContentBin, expectedCode);
      txRespViewModel = HelperTools.JSONDeserialize<SubmitTransactionResponseViewModel>(resp.response.Payload);
      Assert.AreEqual(returnResult, txRespViewModel.ResultDescription);

      // Test submitting multiple  transaction through Binary
      txBytes = HelperTools.HexStringToByteArray(tx6Hex);

      reqContentBin = new ByteArrayContent(txBytes);
      reqContentBin.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Octet);

      resp = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransactions + $"?callbackUrl={url}", client, reqContentBin, expectedCode);
      txsRespViewModel = HelperTools.JSONDeserialize<SubmitTransactionsResponseViewModel>(resp.response.Payload);
      Assert.AreEqual(returnResult, txsRespViewModel.Txs[0].ResultDescription);
    }


    [TestMethod]
    public async Task SubmitTransactionJsonAuthenticated()
    {
      // use special free fee policy for user
      feeQuoteRepositoryMock.FeeFileName = "feeQuotesWithIdentity.json";

      var reqContent = new StringContent($"{{ \"rawtx\": \"{txZeroFeeHex}\" }}");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);

      // txZeroFeeHex - it should fail without authentication
      var response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);
      Assert.AreEqual(0, rpcClientFactoryMock.AllCalls.FilterCalls("mocknode0:sendrawtransactions/").Count()); // no calls, to submit txs since we do not pay enough fee

      var payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();
      // Check if all fields are set
      AssertIsOK(payload, txZeroFeeHash, "failure", "Not enough fees");

      // Test token valid until year 2030. Generate with:
      //    TokenManager.exe generate -n 5 -i http://mysite.com -a http://myaudience.com -k thisisadevelopmentkey -d 3650
      //
      RestAuthentication = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI1IiwibmJmIjoxNTk5NDExNDQzLCJleHAiOjE5MTQ3NzE0NDMsImlhdCI6MTU5OTQxMTQ0MywiaXNzIjoiaHR0cDovL215c2l0ZS5jb20iLCJhdWQiOiJodHRwOi8vbXlhdWRpZW5jZS5jb20ifQ.Z43NASAbIxMZrL2MzbJTJD30hYCxhoAs-8heDjQMnjM";
      // now it should succeed for this user
      response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);

      rpcClientFactoryMock.AllCalls.AssertContains("mocknode0:sendrawtransactions/", "mocknode0:sendrawtransactions/" + txZeroFeeHash);
      payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();

      // Check if all fields are set
      AssertIsOK(payload, txZeroFeeHash);
    }



    [TestMethod]
    public async Task SubmitTransactionJsonFeeQuoteExpired()
    {
      // use special free fee policy 
      feeQuoteRepositoryMock.FeeFileName = "feeQuotesAllFree.json";
      var feeQuotes = feeQuoteRepositoryMock.GetValidFeeQuotesByIdentity(null);
      Assert.AreEqual(1, feeQuotes.Count()); 

      var reqContent = new StringContent($"{{ \"rawtx\": \"{txZeroFeeHex}\" }}");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);

      using (MockedClock.NowIs(DateTime.UtcNow.AddMinutes(FeeQuoteRepositoryMock.quoteExpiryMinutes)))
      {
        feeQuotes = feeQuoteRepositoryMock.GetValidFeeQuotesByIdentity(null);
        Assert.AreEqual(1, feeQuotes.Count()); // should return current

        var response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
        VerifySignature(response);

        rpcClientFactoryMock.AllCalls.AssertContains("mocknode0:sendrawtransactions/", "mocknode0:sendrawtransactions/" + txZeroFeeHash);
        var payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();

        // Check if all fields are set
        AssertIsOK(payload, txZeroFeeHash);
      }

    }

    [TestMethod]
    public async Task SubmitTransactionJsonTwoValidFeeQuotes()
    {
      // use special free fee policy 
      feeQuoteRepositoryMock.FeeFileName = "feeQuotesWithCreatedAt.json";

      var feeQuotes = feeQuoteRepositoryMock.GetValidFeeQuotesByIdentity(null);
      Assert.AreEqual(1, feeQuotes.Count()); // current feeQuote is valid now

      var reqContent = new StringContent($"{{ \"rawtx\": \"{txZeroFeeHex}\" }}");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);

      using (MockedClock.NowIs(new DateTime(2020, 9, 1, 12, 6, 0))) // go back in time
      {
        feeQuotes = feeQuoteRepositoryMock.GetValidFeeQuotesByIdentity(null);
        Assert.AreEqual(2, feeQuotes.Count());

        var response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
        VerifySignature(response);

        rpcClientFactoryMock.AllCalls.AssertContains("mocknode0:sendrawtransactions/", "mocknode0:sendrawtransactions/" + txZeroFeeHash);
        var payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();

        // Check if all fields are set
        AssertIsOK(payload, txZeroFeeHash);
      }

    }

    [TestMethod]
    public async Task SubmitTransactionJsonTwoNodes()
    {
      AddMockNode(1);
      Assert.AreEqual(2, Nodes.GetNodes().Count());


      var reqContent = new StringContent($"{{ \"rawtx\": \"{txC3Hex}\" }}");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);

      var response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);

      rpcClientFactoryMock.AllCalls.AssertContains("mocknode0:sendrawtransactions/", "mocknode0:sendrawtransactions/" + txC3Hash);
      rpcClientFactoryMock.AllCalls.AssertContains("mocknode1:sendrawtransactions/", "mocknode1:sendrawtransactions/" + txC3Hash);
      
      var payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();

      // Check if all fields are set
      AssertIsOK(payload, txC3Hash);
    }

    [TestMethod]
    public async Task SubmitTransactionJsonTwoNodesOneDown()
    {
      AddMockNode(1);
      Assert.AreEqual(2, Nodes.GetNodes().Count());

      rpcClientFactoryMock.DisconnectNode("mocknode0");

      var reqContent = new StringContent($"{{ \"rawtx\": \"{txC3Hex}\" }}");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);

      var response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);

      rpcClientFactoryMock.AllCalls.AssertContains("mocknode1:sendrawtransactions/", "mocknode1:sendrawtransactions/" + txC3Hash);

      var payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();

      // Check if all fields are set
      AssertIsOK(payload, txC3Hash);
    }


    [TestMethod]
    public async Task SubmitTransactionJsonBothNodesDown_should_return_failure()
    {
      AddMockNode(1);
      Assert.AreEqual(2, Nodes.GetNodes().Count());

      rpcClientFactoryMock.DisconnectNode("mocknode0");
      

      var reqContent = new StringContent($"{{ \"rawtx\": \"{txC3Hex}\" }}");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);


      // fetch blockchain info while one node is sitll available and then last connected node
      _ = blockChainInfo.GetInfo();
      rpcClientFactoryMock.DisconnectNode("mocknode1");

      var response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);

      var payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();

      Assert.AreEqual("failure", payload.ReturnResult);
    }


    [TestMethod]
    public async Task SubmitTransactionZeroFeeJson()
    {
      var reqContent = new StringContent($"{{ \"rawtx\": \"{txZeroFeeHex}\" }}");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);

      var response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);

      Assert.AreEqual(0, rpcClientFactoryMock.AllCalls.FilterCalls("mocknode0:sendrawtransactions/").Count()); // no calls, to submit txs since we do not pay enough fee
      
      // We still expect one call to fetch previous outputs (we need to check them to calculate fee) 
      Assert.AreEqual(1, rpcClientFactoryMock.AllCalls.FilterCalls("mocknode0:gettxouts").Count());

      var payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();

      // Check if all fields are set
      AssertIsOK(payload, txZeroFeeHash, "failure", "Not enough fees");
    }

    [TestMethod]
    [OverrideSetting("AppSettings:CheckFeeDisabled", true)]
    public async Task SubmitTransactionJsonCheckFeeDisabled()
    {
      var reqContent = new StringContent($"{{ \"rawtx\": \"{txZeroFeeHex}\" }}");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);

      var response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransaction, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);

      Assert.AreEqual(1, rpcClientFactoryMock.AllCalls.FilterCalls("mocknode0:sendrawtransactions/").Count()); // no calls, to submit txs since we do not pay enough fee

      var payload = response.response.ExtractPayload<SubmitTransactionResponseViewModel>();

      // Check if all fields are set
      AssertIsOK(payload, txZeroFeeHash);      
    }

    void Assert2ValidAnd1Invalid(SubmitTransactionsResponseViewModel response)
    {

      // tx1 and tx2 should be acccepted, bzt txZeroFee should not be
      rpcClientFactoryMock.AllCalls.AssertContains("mocknode0:sendrawtransactions/", "mocknode0:sendrawtransactions/" + txC3Hash + "/" + tx2Hash);

      // validate header
      Assert.AreEqual("1.3.0", response.ApiVersion);
      Assert.IsTrue((DateTime.UtcNow - response.Timestamp).TotalSeconds < 60);
      Assert.AreEqual(MinerId.GetCurrentMinerIdAsync().Result, response.MinerId);
      Assert.AreEqual(BlockChainInfo.GetInfo().BestBlockHeight, response.CurrentHighestBlockHeight);
      Assert.AreEqual(BlockChainInfo.GetInfo().BestBlockHash, response.CurrentHighestBlockHash);

      // validate individual transactions
      Assert.AreEqual(1, response.FailureCount);
      Assert.AreEqual(3, response.Txs.Length);

      // Failures are listed first
      Assert.AreEqual(txZeroFeeHash, response.Txs[0].Txid);
      Assert.AreEqual("failure", response.Txs[0].ReturnResult);
      Assert.AreEqual("Not enough fees", response.Txs[0].ResultDescription);

      Assert.AreEqual(txC3Hash, response.Txs[1].Txid);
      Assert.AreEqual("success", response.Txs[1].ReturnResult);
      Assert.AreEqual("", response.Txs[1].ResultDescription);

      Assert.AreEqual(tx2Hash, response.Txs[2].Txid);
      Assert.AreEqual("success", response.Txs[2].ReturnResult);
      Assert.AreEqual("", response.Txs[2].ResultDescription);
    }

    [TestMethod]
    public async Task SubmitTransactionsJson()
    {
      var reqContent = new StringContent($"[ {{ \"rawtx\": \"{txC3Hex}\" }}, {{ \"rawtx\": \"{txZeroFeeHex}\" }},  {{ \"rawtx\": \"{tx2Hex}\" }}]");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);
      var response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransactions, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);


      Assert2ValidAnd1Invalid(response.response.ExtractPayload<SubmitTransactionsResponseViewModel>());
    }

    [TestMethod]
    public async Task SubmitTransactionsJsonChain()
    {
      // use free fee policy, since the transactions we use are not paying any fee
      feeQuoteRepositoryMock.FeeFileName = "feeQuotesAllFree.json";
      // just register tx0, the rests should get inputs from the bbatc
      rpcClientFactoryMock.SetUpTransaction(txC0Hex);

      var reqContent = new StringContent($"[ {{ \"rawtx\": \"{txC1Hex}\" }}, {{ \"rawtx\": \"{txC2Hex}\" }},  {{ \"rawtx\": \"{txC3Hex}\" }}]");
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Json);
      var response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransactions, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);

      var p = response.response.ExtractPayload<SubmitTransactionsResponseViewModel>();
      Assert.AreEqual(0, p.FailureCount);
      Assert.AreEqual(3,p.Txs.Length);

      Assert.AreEqual(1, rpcClientFactoryMock.AllCalls.FilterCalls("mocknode0:gettxouts").Count(), "just 1 call toward the node was expected");

    }

    [TestMethod]
    public async Task SubmitTransactionsBinary()
    {
      var bytes = HelperTools.HexStringToByteArray(txC3Hex + txZeroFeeHex + tx2Hex); 
      
      var reqContent = new ByteArrayContent(bytes);
      reqContent.Headers.ContentType = new MediaTypeHeaderValue(MediaTypeNames.Application.Octet);
      var response = await Post<SignedPayloadViewModel>(MapiServer.ApiMapiSubmitTransactions, client, reqContent, HttpStatusCode.OK);
      VerifySignature(response);

      Assert2ValidAnd1Invalid(response.response.ExtractPayload<SubmitTransactionsResponseViewModel>());
    }


  }
}
