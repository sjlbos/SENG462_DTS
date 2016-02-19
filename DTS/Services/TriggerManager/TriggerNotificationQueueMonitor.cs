
using System;
using System.Collections.Concurrent;
using System.Globalization;
using log4net;
using Newtonsoft.Json;
using RabbitMessaging;
using ServiceHost;
using TriggerManager.Models;

namespace TriggerManager
{
    /// <summary>
    /// A class that listens for trigger update notifcations sent by the DTS API. 
    /// </summary>
    public class TriggerNotificationQueueMonitor : QueueMonitorWorker
    {
        private static readonly ILog Log = LogManager.GetLogger(typeof (TriggerNotificationQueueMonitor));

        private readonly BlockingCollection<TriggerUpdateNotification> _notificationBuffer; 

        /// <param name="instanceId">The unique name of this worker.</param>
        /// <param name="messageReceiver">The receiver set up to listen for incoming trigger update notifications.</param>
        /// <param name="notificationBuffer">The buffer to which the worker will write all received trigger update notifications.</param>
        public TriggerNotificationQueueMonitor(string instanceId, IMessageReceiver messageReceiver, BlockingCollection<TriggerUpdateNotification> notificationBuffer) : base(instanceId, messageReceiver)
        {
            if (notificationBuffer == null)
                throw new ArgumentNullException("notificationBuffer");

            _notificationBuffer = notificationBuffer;
        }

        /// <summary>
        /// Handles incoming trigger update notifaction messages. When a new message arrives, this method
        /// deserializes the message from JSON into a TriggerUpdateNotification object, which is then placed
        /// on a shared queue.
        /// </summary>
        /// <param name="message">A JSON string representation of a TriggerUpdateNotification.</param>
        public override void ProcessMessage(string message)
        {
            try
            {
                var updateNotification = JsonConvert.DeserializeObject<TriggerUpdateNotification>(message);
                
                Log.DebugFormat(CultureInfo.InvariantCulture,
                    "Worker {0} received an trigger notification message. User: \"{1}\" Trigger Type: \"{2}\"",
                    InstanceId, updateNotification.UserId, updateNotification.TriggerType);

                _notificationBuffer.Add(updateNotification);
            }
            catch (JsonSerializationException ex)
            {
                Log.Error(String.Format(CultureInfo.InvariantCulture,
                    "Worker {0} could not deserialize received trigger update message: \"{1}\"", InstanceId, message), ex);
            }  
        }
    }
}
