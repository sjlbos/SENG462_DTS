
namespace TriggerManager.Models
{
    public class Trigger
    {
        public int Id { get; set; }
        public TriggerType Type { get; set; }
        public int UserId { get; set; }
        public string StockSymbol { get; set; }
        public int NumberOfShares { get; set; }
        public decimal TriggerPrice { get; set; }
    }
}
