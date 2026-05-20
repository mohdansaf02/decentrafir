class UserModel {
  final String id;
  final String name;
  final String email;
  final String role;
  final String? walletAddress;
  final String? department;
  final String? badgeNumber;

  const UserModel({
    required this.id,
    required this.name,
    required this.email,
    required this.role,
    this.walletAddress,
    this.department,
    this.badgeNumber,
  });

  factory UserModel.fromJson(Map<String, dynamic> json) => UserModel(
        id: json['id'] ?? '',
        name: json['name'] ?? '',
        email: json['email'] ?? '',
        role: json['role'] ?? 'citizen',
        walletAddress: json['walletAddress'],
        department: json['department'],
        badgeNumber: json['badgeNumber'],
      );

  bool get isCitizen => role == 'citizen';
  bool get isPolice => role == 'police';
  bool get isAdmin => role == 'admin';
}
