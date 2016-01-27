
using System;
using System.Net;
using System.Text;
using Newtonsoft.Json;

namespace WorkloadGeneratorSlave
{
    internal class ApiCommand 
    {
        public string Id { get; set; }
        public string Method { get; set; }
        public Uri Uri { get; set; }
        public string RequestBody { get; set; }
        public HttpStatusCode ExpectedStatusCode { get; set; }

        private WebRequest _request;

        [JsonIgnore]
        public WebRequest HttpRequest
        {
            get
            {
                if (_request != null)
                    return _request;

                _request = WebRequest.Create(Uri);
                _request.Method = Method;

                if (RequestBody != null)
                {
                    var bodyBytes = Encoding.UTF8.GetBytes(RequestBody);
                    _request.ContentLength = bodyBytes.Length;
                    _request.ContentType = "application/json";
                    using (var requestStream = _request.GetRequestStream())
                    {
                        requestStream.Write(bodyBytes, 0, bodyBytes.Length);
                    }
                }
                return _request;
            }
        }
    }
}
