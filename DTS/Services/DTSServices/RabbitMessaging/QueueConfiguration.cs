
using System.Collections.Generic;

namespace RabbitMessaging
{
    public class QueueConfiguration
    {
        public string Name { get; set; }
        public bool IsDurable { get; set; }
        public bool IsExclusive { get; set; }
        public bool AutoDelete { get; set; }
        public ExchangeConfiguration Exchange { get; set; }
        public IList<string> BindingKeys { get; set; } 
    }
}
