import 'dart:convert';
import 'package:fir_management/core/config/app_config.dart';
import 'package:http/http.dart' as http;

/// Web/MetaMask integration layer.
/// On web, connect via `window.ethereum`; on mobile use WalletConnect in production.
class BlockchainService {
  String? _walletAddress;
  String? get walletAddress => _walletAddress;

  bool get isConnected => _walletAddress != null;

  /// Simulates wallet connect for demo; replace with JS interop on Flutter web.
  Future<String?> connectWallet() async {
    // Production: use dart:js_interop / wallet_connect_flutter_v2
    _walletAddress = '0xDemoWallet${DateTime.now().millisecondsSinceEpoch}';
    return _walletAddress;
  }

  void disconnect() => _walletAddress = null;

  Future<Map<String, dynamic>?> getChainInfo() async {
    try {
      final res = await http.post(
        Uri.parse(AppConfig.rpcUrl),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'jsonrpc': '2.0',
          'method': 'eth_chainId',
          'params': [],
          'id': 1,
        }),
      );
      if (res.statusCode == 200) {
        final data = jsonDecode(res.body);
        return {'chainId': data['result']};
      }
    } catch (_) {}
    return null;
  }

  String explorerTxUrl(String txHash) {
    return 'https://mumbai.polygonscan.com/tx/$txHash';
  }
}
