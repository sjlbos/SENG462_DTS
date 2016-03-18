using System;
using System.Globalization;
using System.Net.Sockets;
using System.Text;
using log4net;

namespace TriggerManager.Repository
{
    public class SocketQuoteProvider : IQuoteProvider
    {
        private static readonly ILog Log = LogManager.GetLogger(typeof (SocketQuoteProvider));

        private const int ResponseBufferSize = 1024;
        private const string QuoteCacheSendFormat = "False,{0},{1}";

        private readonly string _quoteCacheHost;
        private readonly int _quoteCachePort;

        public SocketQuoteProvider(string hostName, int port)
        {
            if(hostName == null)
                throw new ArgumentNullException("hostName");

            _quoteCacheHost = hostName;
            _quoteCachePort = port;
        }

        public decimal GetStockPriceForUser(string stockSymbol, string username)
        {
            if (stockSymbol == null)
                throw new ArgumentNullException("stockSymbol");
            if (username == null)
                throw new ArgumentNullException("username");

            string requestMessage = GetQuoteRequestString(stockSymbol, username);
            string response = GetResponseFromQuoteCache(requestMessage);

            decimal stockPrice;

            try
            {
                stockPrice = Convert.ToDecimal(response);
            }
            catch (FormatException ex)
            {
                throw new QuoteException(String.Format(CultureInfo.InvariantCulture, "Received an invalid quote response: \"{0}\"", response), ex);
            }

            if (stockPrice < 0)
                throw new QuoteException(String.Format(CultureInfo.InvariantCulture, "The quote cache returned and error for stock \"{0}\" and user \"{1}\".", stockSymbol, username));

            return stockPrice;
        }

        private string GetResponseFromQuoteCache(string requestMessage)
        {
            var messageBytes = Encoding.ASCII.GetBytes(requestMessage);

            try
            {
                // Open a new socket
                var sender = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
                sender.Connect(_quoteCacheHost, _quoteCachePort);

                // Send request
                int sentBytes = sender.Send(messageBytes);
                Log.DebugFormat(CultureInfo.InvariantCulture,
                    "Sent {0} bytes over socket to \"{1}\" on port {2}.", sentBytes, _quoteCacheHost, _quoteCachePort);

                // Wait for response
                var responseBuffer = new byte[ResponseBufferSize];
                int receivedBytes = sender.Receive(responseBuffer);
                Log.DebugFormat(CultureInfo.InvariantCulture,
                    "Received {0} bytes over socket from \"{1}\" on port {2}.", receivedBytes, _quoteCacheHost, _quoteCachePort);

                // Return response as string
                return Encoding.ASCII.GetString(responseBuffer);
            }
            catch (Exception ex)
            {
                throw new QuoteException(String.Format(CultureInfo.InvariantCulture,
                    "Encountered an error while trying to connect to the quote cache at \"{0}\" on port {1}.",
                    _quoteCacheHost, _quoteCachePort), ex);
            }
        }

        private static string GetQuoteRequestString(string stockSymbol, string username)
        {
            return String.Format(CultureInfo.InvariantCulture, QuoteCacheSendFormat, username, stockSymbol);
        } 
    }
}
