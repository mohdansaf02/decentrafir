class FIRModel {
  final String id;
  final String firId;
  final String title;
  final String description;
  final String crimeType;
  final String status;
  final String location;
  final String? ipfsHash;
  final String? transactionHash;
  final int? blockNumber;
  final DateTime? createdAt;

  const FIRModel({
    required this.id,
    required this.firId,
    required this.title,
    required this.description,
    required this.crimeType,
    required this.status,
    required this.location,
    this.ipfsHash,
    this.transactionHash,
    this.blockNumber,
    this.createdAt,
  });

  factory FIRModel.fromJson(Map<String, dynamic> json) => FIRModel(
        id: json['id'] ?? '',
        firId: json['firId'] ?? '',
        title: json['title'] ?? '',
        description: json['description'] ?? '',
        crimeType: json['crimeType'] ?? '',
        status: json['status'] ?? 'pending',
        location: json['location'] ?? '',
        ipfsHash: json['ipfsHash'],
        transactionHash: json['transactionHash'],
        blockNumber: json['blockNumber'],
        createdAt: json['createdAt'] != null ? DateTime.tryParse(json['createdAt']) : null,
      );
}
