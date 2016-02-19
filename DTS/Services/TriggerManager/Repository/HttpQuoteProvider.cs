
using System;
using System.Globalization;
using System.Net;

namespace TriggerManager.Repository
{
    public class HttpQuoteProvider : IQuoteProvider
    {
        private const string QuoteApiPathFormat = "/api/{0}/{1}";

        private readonly string _quoteApiUriFormat;

        public HttpQuoteProvider(string quoteApiRoot)
        {
            if (String.IsNullOrWhiteSpace(quoteApiRoot))
                throw new ArgumentException("Parameter \"quoteApiRoot\" cannot be null or empty.");

            _quoteApiUriFormat = quoteApiRoot + QuoteApiPathFormat;
        }

        public decimal GetStockPriceForUser(string stockSymbol, string username)
        {
            if (String.IsNullOrWhiteSpace(stockSymbol))
                throw new ArgumentException("Parameter \"stockSymbol\" cannot be null or empty.");
            if (String.IsNullOrWhiteSpace(username))
                throw new ArgumentException("Parameter \"username\" cannot be null or empty.");

            var request = GetQuoteCacheApiRequest(stockSymbol, username);

            try
            {
                using (var response = (HttpWebResponse) request.GetResponse())
                {
                    if (response.StatusCode != HttpStatusCode.OK)
                    {
                        throw new RepositoryException(String.Format(CultureInfo.InvariantCulture,
                            "Quote cache responded with status code \"{0}\".", response.StatusCode));
                    }

                    return GetPriceFromWebResponse(response);
                }
            }
            catch (WebException ex)
            {
                throw new RepositoryException("Encountered an error while making a request to: " + request.RequestUri, ex);
            }
        }

        private WebRequest GetQuoteCacheApiRequest(string stockSymbol, string username)
        {
            string quoteApiUriString = WebUtility.UrlEncode(String.Format(CultureInfo.InvariantCulture, _quoteApiUriFormat, stockSymbol, username));
            var quoteApiUri = new Uri(quoteApiUriString);
            var request = WebRequest.Create(quoteApiUri);
            request.Method = "GET";
            return request;
        }

        private decimal GetPriceFromWebResponse(HttpWebResponse response)
        {
            throw new NotImplementedException();
        }
    }
}
