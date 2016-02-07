using System;
using System.Collections.Generic;
using System.Globalization;
using System.IO;
using System.Linq;
using System.Xml.Serialization;
using log4net;
using Nancy;
using TransactionEvents;
using TransactionMonitor.Repository;

namespace TransactionMonitor.Api
{
    public class AuditModule : NancyModule
    {
        private static readonly ILog Log = LogManager.GetLogger(typeof (AuditModule));

        private readonly IAuditRepository _repository;

        public AuditModule(IAuditRepository repository) : base("/audit")
        {
            _repository = repository;

            Get["/transactions"] = parameters =>
            {
                var startTimeParam = Request.Query["start"];
                var endTimeParam = Request.Query["end"];

                try
                {
                    var startTime = String.IsNullOrWhiteSpace(startTimeParam)
                        ? DateTime.MinValue
                        : Convert.ToDateTime(startTimeParam);
                    var endTime = String.IsNullOrWhiteSpace(endTimeParam)
                        ? DateTime.MaxValue
                        : Convert.ToDateTime(endTimeParam);

                    IEnumerable<TransactionEvent> queryResults = _repository.GetAllLogs(startTime, endTime);

                    return new Response
                    {
                        StatusCode = HttpStatusCode.OK,
                        ContentType = "application/xml",
                        Contents = stream =>
                        {
                            XmlLog.Write(queryResults, stream);
                            stream.Flush();
                            stream.Close();
                        },
                    };
                }
                catch (FormatException ex)
                {
                    return new Response
                    {
                        StatusCode = HttpStatusCode.BadRequest,
                        ReasonPhrase = ex.Message
                    };
                }
            };

            Get["/transactions/{FileName}"] = parameters =>
            {
                var fileName = (string) parameters.FileName;
                var startTimeParam = Request.Query["start"];
                var endTimeParam = Request.Query["end"];
                return BuildFileDownloadResponse(fileName, "Hello World");
            };

            Get["/transactions/{Id}"] = parameters =>
            {
                var transactionId = parameters.Id;
                return null;
            };

            Get["/transactions/{Id}/{FileName}"] = parameters =>
            {
                var transactionId = parameters.Id;
                var fileName = (string) parameters.FileName;
                return BuildFileDownloadResponse(fileName, "Hello World");
            };

            Get["/users/{Id}"] = parameters =>
            {
                var userId = parameters.Id;
                var startTimeParam = Request.Query["start"];
                var endTimeParam = Request.Query["end"];

                return null;
            };

            Get["/users/{Id}/{FileName}"] = parameters =>
            {
                var userId = (string) parameters.Id;
                var fileName = (string) parameters.FileName;
                var startTimeParam = Request.Query["start"];
                var endTimeParam = Request.Query["end"];

                return BuildFileDownloadResponse(fileName, "Hello World");
            };
        }

        private Response BuildFileDownloadResponse(string fileName, string fileContent)
        {
            var response = new Response();
            response.Headers.Add("Content-Disposition", String.Format(CultureInfo.InvariantCulture,
                "attachment; filename={0}", fileName));
            response.ContentType = "application/xml";
            response.Contents = stream =>
            {
                using (var writer = new StreamWriter(stream))
                {
                    writer.Write(fileContent);
                }
            };
            return response;
        }
    }
}
