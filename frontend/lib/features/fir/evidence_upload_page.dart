import 'package:dio/dio.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:fir_management/core/providers/auth_provider.dart';

class EvidenceUploadPage extends ConsumerStatefulWidget {
  const EvidenceUploadPage({super.key, required this.firId});
  final String firId;

  @override
  ConsumerState<EvidenceUploadPage> createState() => _EvidenceUploadPageState();
}

class _EvidenceUploadPageState extends ConsumerState<EvidenceUploadPage> {
  bool _uploading = false;

  Future<void> _pickAndUpload() async {
    final result = await FilePicker.platform.pickFiles(withData: true);
    if (result == null || result.files.isEmpty) return;
    final file = result.files.first;
    if (file.bytes == null) return;

    setState(() => _uploading = true);
    try {
      final api = ref.read(apiServiceProvider);
      final form = FormData.fromMap({
        'firId': widget.firId,
        'file': MultipartFile.fromBytes(file.bytes!, filename: file.name),
      });
      final res = await api.client.post('/api/v1/evidence/upload', data: form);
      Fluttertoast.showToast(msg: 'Uploaded: ${res.data['ipfsCid']}');
    } catch (_) {
      Fluttertoast.showToast(msg: 'Upload failed');
    } finally {
      if (mounted) setState(() => _uploading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Upload Evidence')),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text('FIR: ${widget.firId}'),
            const SizedBox(height: 24),
            FilledButton.icon(
              onPressed: _uploading ? null : _pickAndUpload,
              icon: const Icon(Icons.cloud_upload),
              label: Text(_uploading ? 'Uploading...' : 'Select File'),
            ),
            const Padding(
              padding: EdgeInsets.all(24),
              child: Text('Supported: PDF, JPG, PNG, MP4, DOC (max 10MB)'),
            ),
          ],
        ),
      ),
    );
  }
}
