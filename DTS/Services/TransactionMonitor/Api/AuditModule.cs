using System;
using System.Collections.Generic;
using System.Configuration;
using System.Globalization;
using System.Linq;
using Nancy;
using TransactionEvents;
using TransactionMonitor.Repository;

namespace TransactionMonitor.Api
{
    public class AuditModule : NancyModule
    {
        private readonly IAuditRepository _repository;
        private readonly string _hostName;

        public AuditModule(IAuditRepository repository) : base("/audit")
        {
            _repository = repository;
            _hostName = ConfigurationManager.AppSettings["ApiRoot"];

            Get["/transactions"] = parameters =>
            {
                var startTimeParam = (string)Request.Query["start"];
                var endTimeParam = (string)Request.Query["end"];
                DateTime startTime;
                DateTime endTime;

                try
                {
                    startTime = String.IsNullOrWhiteSpace(startTimeParam)
                        ? DateTime.MinValue
                        : Convert.ToDateTime(startTimeParam);
                }
                catch (FormatException)
                {
                    return new Response
                    {
                        StatusCode = HttpStatusCode.BadRequest,
                        ReasonPhrase = "Parameter \"start\" is not of a valid format."
                    };
                }

                try
                {
                    endTime = String.IsNullOrWhiteSpace(endTimeParam)
                        ? DateTime.MaxValue
                        : Convert.ToDateTime(endTimeParam);
                }
                catch (FormatException)
                {
                    return new Response
                    {
                        StatusCode = HttpStatusCode.BadRequest,
                        ReasonPhrase = "Parameter \"end\" is not of a valid format."
                    };
                }

                LogDumplogEvent(Request, null);
                var queryResults = _repository.GetAllLogs(startTime, endTime);
                return BuildXmlResponse(queryResults);
            };

            Get["/transactions/{FileName}"] = parameters =>
            {
                var fileName = (string) parameters.FileName;
                var startTimeParam = (string) Request.Query["start"];
                var endTimeParam = (string) Request.Query["end"];
                DateTime startTime;
                DateTime endTime;

                try
                {
                    startTime = String.IsNullOrWhiteSpace(startTimeParam)
                        ? DateTime.MinValue
                        : Convert.ToDateTime(startTimeParam);
                }
                catch (FormatException)
                {
                    return new Response
                    {
                        StatusCode = HttpStatusCode.BadRequest,
                        ReasonPhrase = "Parameter \"start\" is not of a valid format."
                    };
                }

                try
                {
                    endTime = String.IsNullOrWhiteSpace(endTimeParam)
                        ? DateTime.MaxValue
                        : Convert.ToDateTime(endTimeParam);
                }
                catch (FormatException)
                {
                    return new Response
                    {
                        StatusCode = HttpStatusCode.BadRequest,
                        ReasonPhrase = "Parameter \"end\" is not of a valid format."
                    };
                }

                LogDumplogEvent(Request, null);
                var queryResults = _repository.GetAllLogs(startTime, endTime);
                return BuildFileDownloadResponse(fileName, queryResults);
            };

            Get["/users/{Id}"] = parameters =>
            {
                var userId = parameters.Id;
                var startTimeParam = (string) Request.Query["start"];
                var endTimeParam = (string) Request.Query["end"];
                DateTime startTime;
                DateTime endTime;

                try
                {
                    startTime = String.IsNullOrWhiteSpace(startTimeParam)
                        ? DateTime.MinValue
                        : Convert.ToDateTime(startTimeParam);
                }
                catch (FormatException)
                {
                    return new Response
                    {
                        StatusCode = HttpStatusCode.BadRequest,
                        ReasonPhrase = "Parameter \"start\" is not of a valid format."
                    };
                }

                try
                {
                    endTime = String.IsNullOrWhiteSpace(endTimeParam)
                        ? DateTime.MaxValue
                        : Convert.ToDateTime(endTimeParam);
                }
                catch (FormatException)
                {
                    return new Response
                    {
                        StatusCode = HttpStatusCode.BadRequest,
                        ReasonPhrase = "Parameter \"end\" is not of a valid format."
                    };
                }

                LogDumplogEvent(Request, userId);
                var queryResults = _repository.GetLogsForUser(userId, startTime, endTime);
                return BuildXmlResponse(queryResults);
            };

            Get["/users/{Id}/{FileName}"] = parameters =>
            {
                var userId = (string) parameters.Id;
                var fileName = (string) parameters.FileName;
                var startTimeParam = (string) Request.Query["start"];
                var endTimeParam = (string) Request.Query["end"];
                DateTime startTime;
                DateTime endTime;

                try
                {
                    startTime = String.IsNullOrWhiteSpace(startTimeParam)
                        ? DateTime.MinValue
                        : Convert.ToDateTime(startTimeParam);
                }
                catch (FormatException)
                {
                    return new Response
                    {
                        StatusCode = HttpStatusCode.BadRequest,
                        ReasonPhrase = "Parameter \"start\" is not of a valid format."
                    };
                }

                try
                {
                    endTime = String.IsNullOrWhiteSpace(endTimeParam)
                        ? DateTime.MaxValue
                        : Convert.ToDateTime(endTimeParam);
                }
                catch (FormatException)
                {
                    return new Response
                    {
                        StatusCode = HttpStatusCode.BadRequest,
                        ReasonPhrase = "Parameter \"end\" is not of a valid format."
                    };
                }

                LogDumplogEvent(Request, userId);
                var queryResults = _repository.GetLogsForUser(userId, startTime, endTime);
                return BuildFileDownloadResponse(fileName, queryResults);
            };
        }

        private Response BuildXmlResponse(IEnumerable<TransactionEvent> transactionEvents)
        {
            return new Response
            {
                StatusCode = HttpStatusCode.OK,
                ContentType = "application/xml",
                Contents = stream =>
                {
                    XmlLog.Write(transactionEvents, stream);
                    stream.Flush();
                    stream.Close();
                },
            };   
        }

        private Response BuildFileDownloadResponse(string fileName, IEnumerable<TransactionEvent> transactionEvents)
        {
            var response = new Response();
            response.Headers.Add("Content-Disposition", String.Format(CultureInfo.InvariantCulture,
                "attachment; filename={0}", fileName));
            response.ContentType = "application/xml";
            response.Contents = stream =>
            {
                XmlLog.Write(transactionEvents, stream);
                stream.Flush();
                stream.Close();
            };

            return response;
        }

        private void LogDumplogEvent(Request request, string userId)
        {
            string transactionNumberHeader = request.Headers["X-TransNo"].FirstOrDefault();
            if (String.IsNullOrEmpty(transactionNumberHeader))
                return;

            int transactionId = Int32.Parse(transactionNumberHeader);

            var dumplogEvent = new UserCommandEvent
            {
                Id = Guid.NewGuid(),
                TransactionId = transactionId,
                Command = CommandType.DUMPLOG,
                UserId = userId,
                Service = "Transaction Monitor",
                Server = _hostName,
                OccuredAt = DateTime.Now
            };

            _repository.LogUserCommandEvent(dumplogEvent);
        }
    }
}
