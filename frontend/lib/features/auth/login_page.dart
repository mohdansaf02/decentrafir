import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:go_router/go_router.dart';
import 'package:fir_management/core/providers/auth_provider.dart';
import 'package:fir_management/data/services/blockchain_service.dart';

final blockchainServiceProvider = Provider((ref) => BlockchainService());

class LoginPage extends ConsumerStatefulWidget {
  const LoginPage({super.key});

  @override
  ConsumerState<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends ConsumerState<LoginPage> {
  final _email = TextEditingController();
  final _password = TextEditingController();

  @override
  Widget build(BuildContext context) {
    final auth = ref.watch(authProvider);
    final chain = ref.watch(blockchainServiceProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Login')),
      body: Center(
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 420),
          child: Padding(
            padding: const EdgeInsets.all(24),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                TextField(controller: _email, decoration: const InputDecoration(labelText: 'Email'), keyboardType: TextInputType.emailAddress),
                const SizedBox(height: 12),
                TextField(controller: _password, decoration: const InputDecoration(labelText: 'Password'), obscureText: true),
                const SizedBox(height: 24),
                if (auth.loading) const CircularProgressIndicator() else FilledButton(
                  onPressed: () async {
                    final ok = await ref.read(authProvider.notifier).login(_email.text.trim(), _password.text);
                    if (!mounted) return;
                    if (ok) {
                      final role = ref.read(authProvider).user!.role;
                      context.go(role == 'police' ? '/police' : role == 'admin' ? '/admin' : '/citizen');
                    } else {
                      Fluttertoast.showToast(msg: auth.error ?? 'Login failed');
                    }
                  },
                  child: const Text('Login'),
                ),
                const SizedBox(height: 12),
                OutlinedButton(
                  onPressed: () async {
                    await chain.connectWallet();
                    Fluttertoast.showToast(msg: 'Wallet: ${chain.walletAddress}');
                  },
                  child: Text(chain.isConnected ? 'Wallet Connected' : 'Connect MetaMask'),
                ),
                TextButton(onPressed: () => context.go('/register'), child: const Text('Create account')),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
