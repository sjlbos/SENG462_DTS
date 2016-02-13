
namespace TriggerManager.Repository
{
    public interface IQuoteProvider
    {
        decimal GetStockPriceForUser(string stockSymbol, string username);
    }
}
