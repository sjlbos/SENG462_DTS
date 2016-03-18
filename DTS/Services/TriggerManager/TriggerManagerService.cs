
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
        private string _quoteCacheHost;
        private int _quoteCachePort;
        private int _quoteCachePollRateMs;
        private IList<string> _triggerRepoConnectionStrings;

        protected override void InitializeService()
        {
            _dtsApiRoot = ConfigurationManager.AppSettings["DtsApiRoot"];
            
            _quoteCachePollRateMs = Int32.Parse(ConfigurationManager.AppSettings["QuoteCachePollRateMilliseconds"]);
            _quoteCacheHost = ConfigurationManager.AppSettings["QuoteCacheHost"];
            _quoteCachePort = Int32.Parse(ConfigurationManager.AppSettings["QuoteCachePort"]);

            _triggerRepoConnectionStrings = new List<string>();
            foreach (ConnectionStringSettings element in ConfigurationManager.ConnectionStrings)
            {
                _triggerRepoConnectionStrings.Add(element.ConnectionString);
            }
        }

        protected override IList<IWorker> GetWorkerList()
        {
            var notificationReceiver = RabbitMessengerFactory.GetReceiver("TriggerNotificationQueueReceiver");
            var sharedBuffer = new BlockingCollection<TriggerUpdateNotification>();
            var repostiory = new PostgresTriggerRepository(_triggerRepoConnectionStrings);
            var quoteProvider = new SocketQuoteProvider(_quoteCacheHost, _quoteCachePort);
            var controller = new TriggerController(new DtsApiTriggerAuthority(_dtsApiRoot));

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
