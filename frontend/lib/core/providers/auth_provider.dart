import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fir_management/data/models/user_model.dart';
import 'package:fir_management/data/services/api_service.dart';

final apiServiceProvider = Provider((ref) => ApiService());

class AuthState {
  final UserModel? user;
  final bool loading;
  final String? error;

  const AuthState({this.user, this.loading = false, this.error});

  AuthState copyWith({UserModel? user, bool? loading, String? error}) =>
      AuthState(user: user ?? this.user, loading: loading ?? this.loading, error: error);
}

class AuthNotifier extends StateNotifier<AuthState> {
  AuthNotifier(this._api) : super(const AuthState());

  final ApiService _api;

  Future<bool> login(String email, String password) async {
    state = state.copyWith(loading: true, error: null);
    try {
      final res = await _api.client.post('/login', data: {'email': email, 'password': password});
      await _api.saveTokens(res.data['accessToken'], res.data['refreshToken']);
      final user = UserModel.fromJson(res.data['user']);
      state = AuthState(user: user, loading: false);
      return true;
    } catch (e) {
      state = state.copyWith(loading: false, error: 'Login failed');
      return false;
    }
  }

  Future<bool> register(Map<String, dynamic> data) async {
    state = state.copyWith(loading: true, error: null);
    try {
      final res = await _api.client.post('/register', data: data);
      await _api.saveTokens(res.data['accessToken'], res.data['refreshToken']);
      final user = UserModel.fromJson(res.data['user']);
      state = AuthState(user: user, loading: false);
      return true;
    } catch (e) {
      state = state.copyWith(loading: false, error: 'Registration failed');
      return false;
    }
  }

  Future<void> logout() async {
    await _api.clearTokens();
    state = const AuthState();
  }
}

final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  return AuthNotifier(ref.watch(apiServiceProvider));
});
