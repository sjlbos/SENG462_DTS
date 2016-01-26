
using System.Globalization;
using log4net;
using RabbitMessaging;
using ServiceHost;

namespace WorkloadGeneratorSlave
{
    public class WorkloadGeneratorWorker : QueueMonitorWorker
    {
        private static readonly ILog Log = LogManager.GetLogger(typeof (WorkloadGeneratorWorker)); 

        public WorkloadGeneratorWorker(string instanceId, IMessageReceiver messageReceiver) 
            : base(instanceId, messageReceiver)
        {
            
        }

        public override void ProcessMessage(string message)
        {
            Log.DebugFormat(CultureInfo.InvariantCulture,
                "Worker received message: {0}", message);
        }
    }
}
