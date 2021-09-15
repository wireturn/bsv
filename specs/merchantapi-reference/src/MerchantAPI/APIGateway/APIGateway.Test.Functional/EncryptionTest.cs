﻿// Copyright(c) 2020 Bitcoin Association.
// Distributed under the Open BSV software license, see the accompanying file LICENSE

using MerchantAPI.APIGateway.Domain.Actions;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using Sodium;
using System.Security.Cryptography;

namespace MerchantAPI.APIGateway.Test.Functional
{

  [TestClass]
  public class EncryptionTest
  {

    [TestMethod]
    public void TestEncryption()
    {
      var recipientKeyPair = PublicKeyBox.GenerateKeyPair();
      var encryptionKey = MapiEncryption.GetEncryptionKey(recipientKeyPair);
      string s = "Test message";
      var encrypted = MapiEncryption.Encrypt(s, encryptionKey);

      var decrypted = MapiEncryption.Decrypt(encrypted, recipientKeyPair);
      Assert.AreEqual(s, decrypted);

      encrypted[5] ^= 1;

      Assert.ThrowsException<CryptographicException>(() => MapiEncryption.Decrypt(encrypted, recipientKeyPair));

    }
  }
}
