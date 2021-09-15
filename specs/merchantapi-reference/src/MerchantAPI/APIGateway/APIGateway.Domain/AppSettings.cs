﻿// Copyright(c) 2020 Bitcoin Association.
// Distributed under the Open BSV software license, see the accompanying file LICENSE

using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using MerchantAPI.Common.Authentication;
using Microsoft.Extensions.Options;

namespace MerchantAPI.APIGateway.Domain
{
  public class MinerIdServer
  {
    [Required]
    [Url]
    public string Url { get; set; } 
    public string Authentication { get; set; }
    [Required]
    public string Alias { get; set; }
  }

  public class Notification
  {
    [Range(1, int.MaxValue)]
    public int NotificationIntervalSec { get; set; } = 60;

    [Range(2, 100)]
    public int InstantNotificationsTasks { get; set; }

    [Required]
    public int? InstantNotificationsQueueSize { get; set; }

    [Required]
    public int? MaxNotificationsInBatch { get; set; }

    [Required]
    public int? SlowHostThresholdInMs { get; set; }

    [Range(1, 100)]
    public int InstantNotificationsSlowTaskPercentage { get; set; }

    [Required]
    public int? NoOfSavedExecutionTimes { get; set; }

    [Required]
    public int? NotificationsRetryCount { get; set; }

    [Required]
    public int? SlowHostResponseTimeoutMS { get; set; }

    [Required]
    public int? FastHostResponseTimeoutMS { get; set; }
  }

  public class AppSettings : CommonAppSettings
  {
    [Range(1, double.MaxValue)]
    public double QuoteExpiryMinutes { get; set; } = 10;
    public string WifPrivateKey { get; set; }
    
    public MinerIdServer MinerIdServer { get; set; }

    public int MaxBlockChainLengthForFork { get; set; } = 288;

    [Range(1, int.MaxValue)]
    public int ZmqConnectionTestIntervalSec { get; set; } = 60;
 
    public int DeltaBlockHeightForDoubleSpendCheck { get; set; } = 144;

    [Range(1, int.MaxValue)]
    public int CleanUpTxAfterDays { get; set; } = 3;

    [Range(600, int.MaxValue)]
    public int CleanUpTxPeriodSec { get; set; } = 3600;

    [Range(1, int.MaxValue)]
    public int DSHostBanTimeSec { get; set; }

    [Range(2, int.MaxValue)]
    public int DSMaxNumOfTxQueries { get; set; }

    [Range(1, int.MaxValue)]
    public int DSCachedTxRequestsCooldownPeriodSec { get; set; }

    [Range(1, int.MaxValue)]
    public int DSMaxNumOfUnknownTxQueries { get; set; }

    [Range(1, int.MaxValue)]
    public int DSUnknownTxQueryCooldownPeriodSec { get; set; }

    [Range(1, int.MaxValue)]
    public int DSScriptValidationTimeoutSec { get; set; }

    public Notification Notification { get; set; }

    public bool CheckFeeDisabled { get; set; } = false;

    public bool EnableHTTP { get; set; } = false;
  }

  public class AppSettingValidator : IValidateOptions<AppSettings>
  {
    public ValidateOptionsResult Validate(string name, AppSettings options)
    {
      var hasURL = options.MinerIdServer != null && !string.IsNullOrEmpty(options.MinerIdServer.Url);
      if (!hasURL && string.IsNullOrWhiteSpace(options.WifPrivateKey)
      ||
       hasURL && !string.IsNullOrWhiteSpace(options.WifPrivateKey)
      )
      {
        return ValidateOptionsResult.Fail(
          $"Invalid configuration - either  {nameof(AppSettings.MinerIdServer)} or {nameof(AppSettings.WifPrivateKey)} must be specified.");
      }

      // Explicitly trigger validation of nested objects
      if (options.MinerIdServer != null && !string.IsNullOrEmpty(options.MinerIdServer.Url))
      {
        var validationResults = new List<ValidationResult>();
        var validationContext = new ValidationContext(options.MinerIdServer, serviceProvider: null, items: null);
        if (!Validator.TryValidateObject(options.MinerIdServer, validationContext, validationResults, true))
        {
          return ValidateOptionsResult.Fail(string.Join(",", validationResults.Select(x => x.ErrorMessage).ToArray()));
        }
      }

      if (options.Notification == null)
      {
        return ValidateOptionsResult.Fail(
          $"Invalid configuration -  {nameof(AppSettings.Notification)} settings must be specified.");
      }
      else
      {
        var validationResults = new List<ValidationResult>();
        var validationContext = new ValidationContext(options.Notification, serviceProvider: null, items: null);
        if (!Validator.TryValidateObject(options.Notification, validationContext, validationResults, true))
        {
          return ValidateOptionsResult.Fail(string.Join(",", validationResults.Select(x => x.ErrorMessage).ToArray()));
        }
      }

      return ValidateOptionsResult.Success;
    }
  }
}
