using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Globalization;
using System.Threading;
using log4net;
using ServiceHost;
using TriggerManager.Models;
using TriggerManager.Repository;

namespace TriggerManager
{
    /// <summary>
    /// A class responsible for reading the TriggerUpdateNotification objects received by the TriggerNotificationQueueMonitor, retrieving the
    /// correct triggers from the DTS repository in response to received notifications, and updating the Trigger Monitor Service's list of 
    /// active triggers.
    /// </summary>
    public class RepositoryReaderWorker : IWorker
    {
        private static readonly ILog Log = LogManager.GetLogger(typeof (RepositoryReaderWorker));

        public string InstanceId { get; private set; }
        private readonly BlockingCollection<TriggerUpdateNotification> _updateBuffer; 
        private readonly ITriggerRepository _repository;
        private readonly TriggerController _controller;

        private readonly IDictionary<int, DateTime> _lastUpdated;

        /// <param name="instanceId">The unique name of the worker instance.</param>
        /// <param name="updateBuffer">The buffer from which to read TriggerUpdateNotification messages.</param>
        /// <param name="repository">The DTS trigger repository.</param>
        /// <param name="controller">The TriggerController instance where triggers will be updated.</param>
        public RepositoryReaderWorker(string instanceId, BlockingCollection<TriggerUpdateNotification> updateBuffer, ITriggerRepository repository, TriggerController controller)
        {
            if (instanceId == null)
                throw new ArgumentNullException("instanceId");
            if (updateBuffer == null)
                throw new ArgumentNullException("updateBuffer");
            if (repository == null)
                throw new ArgumentNullException("repository");
            if (controller == null)
                throw new ArgumentNullException("controller");

            InstanceId = instanceId;
            _updateBuffer = updateBuffer;
            _repository = repository;
            _controller = controller;
            _lastUpdated = new Dictionary<int, DateTime>();
        }

        /// <summary>
        /// Continually monitors the RepositoryReaderWorker's TriggerUpdateNotification buffer. When a new 
        /// notification arrives, the worker retrives the correct trigger(s) from the DTS repository and updates
        /// its controller's list of active triggers.
        /// </summary>
        /// <param name="cancellationToken">A cancellation token used to stop the worker.</param>
        public void Run(CancellationToken cancellationToken)
        {
            if (cancellationToken.IsCancellationRequested)
                return;

            // Get initial list of triggers
            Log.Info("Getting initial trigger list from the repostiory...");
            var triggerList = _repository.GetAllTriggers();
            _controller.UpdateTriggers(triggerList);
            Log.Info("Triggers were successfully retrieved.");

            while (true)
            {
                if (cancellationToken.IsCancellationRequested)
                    return;

                try
                {
                    var nextNotification = _updateBuffer.Take(cancellationToken);
                    HandleTriggerNotification(nextNotification);
                }
                catch (OperationCanceledException)
                {
                    return;
                }
            }
        }

        private void HandleTriggerNotification(TriggerUpdateNotification updateNotification)
        {
            if (updateNotification == null)
            {
                Log.WarnFormat(CultureInfo.InvariantCulture, "Worker {0} received a null TriggerUpdateNotification.", InstanceId);
                return;
            }

            if (!NotificationForcesUpdate(updateNotification))
            {
                Log.DebugFormat(CultureInfo.InvariantCulture,
                    "Skipping trigger update for User={0}, TransactionId={1}. No update required.", 
                    updateNotification.UserId, updateNotification.TransactionId);
                return;
            }

            UpdateUserTriggers(updateNotification);
        }

        private bool NotificationForcesUpdate(TriggerUpdateNotification updateNotification)
        {
            // Only update the trigger list if the trigger is for an unknown user or if the user's triggers have been
            // updated since the last pull from the trigger repository. This reduces repository load if a users makes quick,
            // successive updates to their trigger list.
            return !_lastUpdated.ContainsKey(updateNotification.UserId) ||
                   updateNotification.UpdatedAt >= _lastUpdated[updateNotification.UserId];
        }

        private void UpdateUserTriggers(TriggerUpdateNotification updateNotification)
        {
            if (updateNotification.TriggerType == TriggerType.Buy)
            {
                Log.DebugFormat(CultureInfo.InvariantCulture, "Updating buy triggers for user with Id={0}.", updateNotification.UserId);
                UpdateBuyTriggersForUser(updateNotification.UserId);
            }
            else
            {
                Log.DebugFormat(CultureInfo.InvariantCulture, "Updating sell triggers for user with Id={0}.", updateNotification.UserId);
                UpdateSellTriggersForUser(updateNotification.UserId);
            }
        }

        private void UpdateBuyTriggersForUser(int userId)
        {
            UpdateLastUpdatedListForUser(userId, DateTime.UtcNow);
            var buyTriggers = _repository.GetBuyTriggersForUser(userId);
            _controller.UpdateBuyTriggersForUser(userId, buyTriggers);
        }

        private void UpdateSellTriggersForUser(int userId)
        {
            UpdateLastUpdatedListForUser(userId, DateTime.UtcNow);
            var sellTriggers = _repository.GetSellTriggersForUser(userId);
            _controller.UpdateSellTriggersForUser(userId, sellTriggers);
        }

        private void UpdateLastUpdatedListForUser(int userId, DateTime dateTime)
        {
            if (!_lastUpdated.ContainsKey(userId))
            {
                _lastUpdated.Add(userId, dateTime);
            }
            else
            {
                _lastUpdated[userId] = dateTime;
            }
        }
    }
}
