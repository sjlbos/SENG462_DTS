
namespace TriggerManager.Models
{
    public class Trigger
    {
        public int Id { get; set; }
        public TriggerType Type { get; set; }
        public string UserId { get; set; }
        public int UserDbId { get; set; }
        public string StockSymbol { get; set; }
        public int NumberOfShares { get; set; }
        public decimal TriggerPrice { get; set; }
    }
}
