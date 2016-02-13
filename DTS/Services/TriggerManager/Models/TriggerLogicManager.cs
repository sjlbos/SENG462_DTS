
using System.Collections.Generic;
using System.Linq;

namespace TriggerManager.Models
{
    /// <summary>
    /// A thread-safe class for managing when buy and sell triggers should be fired. For each stock, all buy and sell triggers for that stock
    /// are stored in a pair of  lists sorted by their trigger price. This allows for quick evaluation of which triggers will fire 
    /// given a change in stock price.
    /// </summary>
    internal class TriggerLogicManager
    {
        private readonly IDictionary<string, SortedList<decimal, IList<Trigger>>> _buyTriggers;
        private readonly IDictionary<string, SortedList<decimal, IList<Trigger>>> _sellTriggers; 
        private readonly object _syncLock;

        public TriggerLogicManager()
        {
            _buyTriggers = new Dictionary<string, SortedList<decimal, IList<Trigger>>>();
            _sellTriggers = new Dictionary<string, SortedList<decimal, IList<Trigger>>>();
            _syncLock = new object();
        }

        /// <summary>
        /// The complete list of stock symbols the TriggerLogicManager is currently monitoring.
        /// </summary>
        public IEnumerable<string> StockList
        {
            get { return _buyTriggers.Keys.Union(_sellTriggers.Keys); }
        } 

        /// <summary>
        /// Returns a list of buy trigger objects that should be fired if the supplied stock reaches the supplied price.
        /// Buy triggers are fired when the stock price is equal to or less than the trigger price.
        /// </summary>
        /// <param name="stockSymbol">A 3 character stock symbol.</param>
        /// <param name="stockPrice">The current price of the stock corresponding to the supplied stock symbol.</param>
        /// <returns>A list of triggers to fire.</returns>
        public IList<Trigger> GetFiredBuyTriggersAtStockPrice(string stockSymbol, decimal stockPrice)
        {
            var firedTriggers = new List<Trigger>();
            lock (_syncLock)
            {
                if (!_buyTriggers.ContainsKey(stockSymbol))
                    return firedTriggers;

                var priceTriggerMap = _buyTriggers[stockSymbol];

                // Trigger prices are stored in ascending order. We reverse this order to traverse prices 
                // in descending order, firing triggers whos trigger price greater than or equal to the current stock price. 
                foreach (decimal triggerPrice in priceTriggerMap.Keys.Reverse())
                {
                    if (stockPrice <= triggerPrice)
                    {
                        firedTriggers.AddRange(priceTriggerMap[stockPrice]);
                    }
                    else
                    {
                        break;
                    }
                }
                return firedTriggers;
            }
        }

        /// <summary>
        /// Returns a list of sell trigger objects that should be fired if the supplied stock reaches the supplied price.
        /// Sell triggers are fired when the stock price is equal or greater than the trigger price.
        /// </summary>
        /// <param name="stockSymbol">A 3 character stock symbol.</param>
        /// <param name="stockPrice">The current price of the stock corresponding to the supplied stock symbol.</param>
        /// <returns>A list of triggers to fire.</returns>
        public IList<Trigger> GetFiredSellTriggersAtStockPrice(string stockSymbol, decimal stockPrice)
        {
            var firedTriggers = new List<Trigger>();
            lock (_syncLock)
            {
                if (!_sellTriggers.ContainsKey(stockSymbol))
                    return firedTriggers;

                var priceTriggerMap = _sellTriggers[stockSymbol];

                foreach (decimal triggerPrice in priceTriggerMap.Keys)
                {
                    if (stockPrice >= triggerPrice)
                    {
                        firedTriggers.AddRange(priceTriggerMap[triggerPrice]);
                    }
                    else
                    {
                        break;
                    }
                }
                return firedTriggers;
            }
        }

        public void AddTrigger(Trigger trigger)
        {
            lock (_syncLock)
            {
                
            }
        }

        public void RemoveTrigger(Trigger trigger)
        {
            lock (_syncLock)
            {
                
            }
        }

        public void Reset()
        {
            lock (_syncLock)
            {
                
            }
        }
    }
}
