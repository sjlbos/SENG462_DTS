﻿
using System;
using Newtonsoft.Json;
using Newtonsoft.Json.Converters;

namespace TriggerManager.Models
{
    public class TriggerUpdateNotification
    {
        [JsonConverter(typeof(StringEnumConverter))]
        public TriggerType TriggerType { get; set; }
        public int UserDbId { get; set; }
        public string UserId { get; set; }
        public int TransactionId { get; set; }
        public DateTime UpdatedAt { get; set; }
    }
}
