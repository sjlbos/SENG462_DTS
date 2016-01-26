
namespace RabbitMessaging
{
    public class ExchangeConfiguration
    {
        public string Name { get; set; }
        public string ExchangeType { get; set; }
        public bool IsDurable { get; set; }
        public bool AutoDelete { get; set; }
    }
}
