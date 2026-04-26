import grpc
from concurrent import futures
from app.grpc import engine_pb2, engine_pb2_grpc

class EngineServicer(engine_pb2_grpc.EngineServiceServicer):
    def AnalyzeVideo(self, request, context):
        print(f"📥 Received analysis request for video: {request.video_id}")
        
        annotation = engine_pb2.AnnotationEvent(
            timestamp_ms=5000,
            type="comment",
            severity="info",
            comment="gRPC connection working correctly",
            translation="",
            shape=engine_pb2.AnnotationShape(
                shape_type="rect",
                x=0.1,
                y=0.1,
                width=0.2,
                height=0.1,
                color="#FF0000"
            )
        )

        return engine_pb2.AnalyzeVideoResponse(annotations=[annotation])

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    engine_pb2_grpc.add_EngineServiceServicer_to_server(EngineServicer(), server)
    server.add_insecure_port("[::]:50051")
    server.start()
    print("✅ gRPC server running on port 50051")
    server.wait_for_termination()