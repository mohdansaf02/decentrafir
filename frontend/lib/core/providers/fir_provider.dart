import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fir_management/data/models/fir_model.dart';
import 'package:fir_management/core/providers/auth_provider.dart';

final firListProvider = FutureProvider.autoDispose<List<FIRModel>>((ref) async {
  final api = ref.watch(apiServiceProvider);
  final res = await api.client.get('/api/v1/fir/all');
  final list = (res.data['data'] as List).map((e) => FIRModel.fromJson(e)).toList();
  return list;
});

final firAnalyticsProvider = FutureProvider.autoDispose<Map<String, dynamic>>((ref) async {
  final api = ref.watch(apiServiceProvider);
  final res = await api.client.get('/api/v1/fir/analytics');
  return Map<String, dynamic>.from(res.data);
});

final firDetailProvider = FutureProvider.family<FIRModel, String>((ref, id) async {
  final api = ref.watch(apiServiceProvider);
  final res = await api.client.get('/api/v1/fir/$id');
  return FIRModel.fromJson(res.data);
});
