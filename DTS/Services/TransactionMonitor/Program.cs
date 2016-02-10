
using System;
using System.Configuration;
using log4net;
using ServiceHost;

namespace TransactionMonitor
{
    public class Program
    {
        static void Main()
        {
            ILog log = LogManager.GetLogger(typeof(Program));
            try
            {
                string serviceName = ConfigurationManager.AppSettings["ServiceName"];
                string displayName = ConfigurationManager.AppSettings["DisplayName"];
                string serviceDescription = ConfigurationManager.AppSettings["ServiceDescription"];
                ServiceRunner.RunService<TransactionMonitorService>(serviceName, displayName, serviceDescription);
            }
            catch (Exception ex)
            {
                log.Fatal("The Transaction Monitor Service encountered a fatal error.", ex);
            }
        }
    }
}
