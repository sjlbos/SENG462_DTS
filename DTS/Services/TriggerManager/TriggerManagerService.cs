
using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Configuration;
using RabbitMessaging;
using ServiceHost;
using TriggerManager.Models;
using TriggerManager.Repository;

namespace TriggerManager
{
    public class TriggerManagerService : WorkerHost
    {
        private string _dtsApiRoot;
        private string _quoteCacheApiRoot;
        private int _quoteCachePollRateMs;
        private string _triggerRepoConnectionString;

        protected override void InitializeService()
        {
            _dtsApiRoot = ConfigurationManager.AppSettings["DtsApiRoot"];
            _quoteCacheApiRoot = ConfigurationManager.AppSettings["QuoteCacheApiRoot"];
            _quoteCachePollRateMs = Int32.Parse(ConfigurationManager.AppSettings["QuoteCachePollRateMilliseconds"]);
            _triggerRepoConnectionString = ConfigurationManager.ConnectionStrings["DtsDb"].ConnectionString;
        }

        protected override IList<IWorker> GetWorkerList()
        {
            var notificationReceiver = RabbitMessengerFactory.GetReceiver("TriggerNotificationQueueReceiver");
            var sharedBuffer = new BlockingCollection<TriggerUpdateNotification>();
            var repostiory = new PostgresTriggerRepository(_triggerRepoConnectionString);
            var quoteProvider = new HttpQuoteProvider(_quoteCacheApiRoot);
            var controller = new TriggerController(_dtsApiRoot);

            var workerList = new List<IWorker>
            {
                new TriggerNotificationQueueMonitor("QueueMonitor", notificationReceiver, sharedBuffer),
                new RepositoryReaderWorker("RepositoryReader", sharedBuffer, repostiory, controller),
                new StockPricePollingWorker("PollingWorker", _quoteCachePollRateMs, controller, quoteProvider)
            };

            return workerList;
        }
    }
}
