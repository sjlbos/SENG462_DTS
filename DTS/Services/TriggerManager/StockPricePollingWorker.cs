
using System;
using System.Globalization;
using System.Threading;
using log4net;
using ServiceHost;
using TriggerManager.Models;
using TriggerManager.Repository;

namespace TriggerManager
{
    public class StockPricePollingWorker : PollingWorker
    {
        public const string TriggerManagerUsername = "TriggerManager";

        private static readonly ILog Log = LogManager.GetLogger(typeof (StockPricePollingWorker));

        private readonly TriggerController _controller;
        private readonly IQuoteProvider _quoteProvider;

        public StockPricePollingWorker(string id, int pollRateMilliseconds, TriggerController controller, IQuoteProvider quoteProvider) : base(id, pollRateMilliseconds)
        {
            if (controller == null)
                throw new ArgumentNullException("controller");
            if (quoteProvider == null)
                throw new ArgumentNullException("quoteProvider");

            _controller = controller;
            _quoteProvider = quoteProvider;
        }

        protected override void HandleShutdownEvent()
        {
            // No cleanup needed
        }

        protected override void DoWork()
        {
            Log.Info("Polling latest stock prices...");
            //foreach (var stock in _controller.StockList)
            //{
            //    Log.DebugFormat(CultureInfo.InvariantCulture,
            //        "Getting quote for stock \"{0}\"...", stock);
            //    decimal stockPrice = _quoteProvider.GetStockPriceForUser(stock, TriggerManagerUsername);
            //    Log.DebugFormat("Price for stock \"{0}\" was returned as \"{1}\"", stock, stockPrice.ToString("C"));

            //    _controller.HandleStockPriceUpdate(stock, stockPrice);
            //}
            Log.Info("Stock price polling complete.");
        }
    }
}
