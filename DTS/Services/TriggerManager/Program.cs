
using System;
using System.Configuration;
using log4net;
using ServiceHost;

namespace TriggerManager
{
    public class Program
    {
        public static void Main(string[] args)
        {
            ILog log = LogManager.GetLogger(typeof(Program));
            try
            {
                string serviceName = ConfigurationManager.AppSettings["ServiceName"];
                string displayName = ConfigurationManager.AppSettings["DisplayName"];
                string serviceDescription = ConfigurationManager.AppSettings["ServiceDescription"];
                ServiceRunner.RunService<TriggerManagerService>(serviceName, displayName, serviceDescription);
            }
            catch (Exception ex)
            {
                log.Fatal("The Trigger Service encountered a fatal error.", ex);
            }
        }
    }
}
