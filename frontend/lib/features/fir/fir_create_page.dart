import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:go_router/go_router.dart';
import 'package:fir_management/core/providers/auth_provider.dart';

class FIRCreatePage extends ConsumerStatefulWidget {
  const FIRCreatePage({super.key});

  @override
  ConsumerState<FIRCreatePage> createState() => _FIRCreatePageState();
}

class _FIRCreatePageState extends ConsumerState<FIRCreatePage> {
  final _title = TextEditingController();
  final _description = TextEditingController();
  final _crimeType = TextEditingController();
  final _location = TextEditingController();
  bool _loading = false;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Register FIR')),
      body: ListView(
        padding: const EdgeInsets.all(24),
        children: [
          TextField(controller: _title, decoration: const InputDecoration(labelText: 'Title')),
          const SizedBox(height: 12),
          TextField(controller: _description, decoration: const InputDecoration(labelText: 'Description'), maxLines: 4),
          const SizedBox(height: 12),
          TextField(controller: _crimeType, decoration: const InputDecoration(labelText: 'Crime Type')),
          const SizedBox(height: 12),
          TextField(controller: _location, decoration: const InputDecoration(labelText: 'Location')),
          const SizedBox(height: 24),
          FilledButton(
            onPressed: _loading ? null : () async {
              setState(() => _loading = true);
              try {
                final api = ref.read(apiServiceProvider);
                final res = await api.client.post('/api/v1/fir/create', data: {
                  'title': _title.text.trim(),
                  'description': _description.text.trim(),
                  'crimeType': _crimeType.text.trim(),
                  'location': _location.text.trim(),
                });
                if (!mounted) return;
                Fluttertoast.showToast(msg: 'FIR created: ${res.data['firId']}');
                context.go('/fir/${res.data['firId']}');
              } catch (e) {
                Fluttertoast.showToast(msg: 'Failed to create FIR');
              } finally {
                if (mounted) setState(() => _loading = false);
              }
            },
            child: _loading ? const SizedBox(height: 20, width: 20, child: CircularProgressIndicator(strokeWidth: 2)) : const Text('Submit FIR'),
          ),
        ],
      ),
    );
  }
}
