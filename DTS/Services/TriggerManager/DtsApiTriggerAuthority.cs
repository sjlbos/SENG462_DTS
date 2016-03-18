using System;
using System.Globalization;
using System.Net;
using log4net;
using TriggerManager.Models;

namespace TriggerManager
{
    public class DtsApiTriggerAuthority : ITriggerAuthority
    {
        private static readonly ILog Log = LogManager.GetLogger(typeof (DtsApiTriggerAuthority));

        private const string CommitBuyTriggerRoute = "/api/users/{0}/buy-triggers/{1}";
        private const string CommitSellTriggerRoute = "/api/users/{0}/sell-triggers/{1}";

        private readonly Uri _dtsApiRoot;

        public DtsApiTriggerAuthority(string dtsApiRoot)
        {
            if (dtsApiRoot == null)
                throw new ArgumentNullException("dtsApiRoot");

            _dtsApiRoot = new Uri(dtsApiRoot);
        }

        public void ExecuteTrigger(Trigger trigger)
        {
            var request = GetRequestForTrigger(trigger);

            Log.DebugFormat(CultureInfo.InvariantCulture, "Making HTTP request to \"{0}\"...", request);
            try
            {
                using (var response = (HttpWebResponse)request.GetResponse())
                {
                    if (response.StatusCode != HttpStatusCode.OK)
                    {
                        throw new TriggerCommitException("Request completed with status code: " + response.StatusCode);
                    }
                }
            }
            catch (WebException ex)
            {
                throw new TriggerCommitException(String.Format(CultureInfo.InvariantCulture,
                     "Encountered an error while trying to commit trigger with ID={0} to API at \"{1}\".", trigger.Id, request), ex);
            }
            Log.Debug("Request completed successfully.");
        }

        private Uri GetRequestUriFromTrigger(Trigger trigger)
        {
            string endpointString = (trigger.Type == TriggerType.Buy)
                ? String.Format(CultureInfo.InvariantCulture, CommitBuyTriggerRoute, Uri.EscapeDataString(trigger.UserId), trigger.Id)
                : String.Format(CultureInfo.InvariantCulture, CommitSellTriggerRoute, Uri.EscapeDataString(trigger.UserId), trigger.Id);
            return new Uri(_dtsApiRoot, endpointString);
        }

        private WebRequest GetRequestForTrigger(Trigger trigger)
        {
            var request = WebRequest.Create(GetRequestUriFromTrigger(trigger));
            request.Method = "POST";
            return request;
        }
    }
}
