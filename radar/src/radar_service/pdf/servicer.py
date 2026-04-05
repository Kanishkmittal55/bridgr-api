"""PdfExtractionService gRPC servicer implementation."""

import logging

import pymupdf
from radar.services.pdf.v1 import pdf_pb2
from radar.services.pdf.v1 import pdf_pb2_grpc

logger = logging.getLogger(__name__)


class PdfExtractionServicer(pdf_pb2_grpc.PdfExtractionServiceServicer):
    """Implements PdfExtractionService RPCs using PyMuPDF for reliable text extraction."""

    def ExtractText(self, request, context):
        """Extract plain text from PDF bytes."""
        if not request.pdf_content:
            return pdf_pb2.ExtractTextResponse(extracted_text="", success=False)

        try:
            doc = pymupdf.open(stream=request.pdf_content, filetype="pdf")
            text_parts = []
            for page in doc:
                text_parts.append(page.get_text())
            doc.close()
            extracted = "".join(text_parts).strip()
            return pdf_pb2.ExtractTextResponse(
                extracted_text=extracted,
                success=bool(extracted),
            )
        except Exception as e:
            logger.warning("PDF extraction failed: %s", e)
            return pdf_pb2.ExtractTextResponse(extracted_text="", success=False)
