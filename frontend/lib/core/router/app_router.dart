import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:fir_management/core/providers/auth_provider.dart';
import 'package:fir_management/features/auth/login_page.dart';
import 'package:fir_management/features/auth/register_page.dart';
import 'package:fir_management/features/home/home_page.dart';
import 'package:fir_management/features/fir/fir_create_page.dart';
import 'package:fir_management/features/fir/fir_detail_page.dart';
import 'package:fir_management/features/fir/evidence_upload_page.dart';
import 'package:fir_management/features/dashboard/citizen_dashboard.dart';
import 'package:fir_management/features/dashboard/police_dashboard.dart';
import 'package:fir_management/features/dashboard/admin_dashboard.dart';

final appRouterProvider = Provider<GoRouter>((ref) {
  final auth = ref.watch(authProvider);

  return GoRouter(
    initialLocation: '/',
    redirect: (context, state) {
      final loggedIn = auth.user != null;
      final isAuthRoute = state.matchedLocation == '/login' || state.matchedLocation == '/register';
      if (!loggedIn && !isAuthRoute && state.matchedLocation != '/') {
        return '/login';
      }
      if (loggedIn && isAuthRoute) {
        return _dashboardForRole(auth.user!.role);
      }
      return null;
    },
    routes: [
      GoRoute(path: '/', builder: (_, __) => const HomePage()),
      GoRoute(path: '/login', builder: (_, __) => const LoginPage()),
      GoRoute(path: '/register', builder: (_, __) => const RegisterPage()),
      GoRoute(path: '/citizen', builder: (_, __) => const CitizenDashboard()),
      GoRoute(path: '/police', builder: (_, __) => const PoliceDashboard()),
      GoRoute(path: '/admin', builder: (_, __) => const AdminDashboard()),
      GoRoute(path: '/fir/create', builder: (_, __) => const FIRCreatePage()),
      GoRoute(
        path: '/fir/:id',
        builder: (_, state) => FIRDetailPage(id: state.pathParameters['id']!),
      ),
      GoRoute(
        path: '/fir/:id/evidence',
        builder: (_, state) => EvidenceUploadPage(firId: state.pathParameters['id']!),
      ),
    ],
  );
});

String _dashboardForRole(String role) {
  switch (role) {
    case 'police':
      return '/police';
    case 'admin':
      return '/admin';
    default:
      return '/citizen';
  }
}
