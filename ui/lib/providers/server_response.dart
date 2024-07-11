class ServerResponse {
  final int code;
  final String message;
  final dynamic data;
  ServerResponse(this.code, this.message, this.data);

  ServerResponse.fromJson(Map<String, dynamic> json)
      : code = json["code"] as int,
        message = json["message"] as String,
        data = json["data"];
}
