using System;
using System.Collections.Generic;
using System.Globalization;
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

                var queryResults = _repository.GetAllLogs(startTime, endTime);
                return BuildXmlResponse(queryResults);
            };

            Get["/transactions/{FileName}"] = parameters =>
            {
                var fileName = (string) parameters.FileName;
                var startTimeParam = Request.Query["start"];
                var endTimeParam = Request.Query["end"];
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

                var queryResults = _repository.GetAllLogs(startTime, endTime);
                return BuildFileDownloadResponse(fileName, queryResults);
            };

            Get["/users/{Id}"] = parameters =>
            {
                var userId = parameters.Id;
                var startTimeParam = Request.Query["start"];
                var endTimeParam = Request.Query["end"];
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

                var queryResults = _repository.GetLogsForUser(userId, startTime, endTime);
                return BuildXmlResponse(queryResults);
            };

            Get["/users/{Id}/{FileName}"] = parameters =>
            {
                var userId = (string) parameters.Id;
                var fileName = (string) parameters.FileName;
                var startTimeParam = Request.Query["start"];
                var endTimeParam = Request.Query["end"];
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
    }
}
