import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:intl/intl.dart';
import 'package:intl/date_symbol_data_local.dart';

// PASSO 1: main corrigido com a inicialização
void main() {
  initializeDateFormatting('pt_BR', null).then((_) {
    runApp(const MyApp());
  });
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Gerenciamento de Ponto',
      theme: ThemeData(
        fontFamily: 'Inter',
        brightness: Brightness.dark,
        scaffoldBackgroundColor: const Color(0xFF0d1117),
        primaryColor: const Color(0xFF58a6ff),
        colorScheme: const ColorScheme.dark(
          primary: Color(0xFF58a6ff),
          secondary: Color(0xFF238636),
          background: Color(0xFF0d1117),
          surface: Color(0xFF161b22),
          onPrimary: Colors.white,
          onSecondary: Colors.white,
          onBackground: Color(0xFFc9d1d9),
          onSurface: Color(0xFFc9d1d9),
          error: Colors.red,
          onError: Colors.white,
        ),
        textTheme: const TextTheme(
          bodyMedium: TextStyle(color: Color(0xFFc9d1d9)),
          headlineSmall: TextStyle(color: Color(0xFFc9d1d9), fontWeight: FontWeight.w600),
          titleLarge: TextStyle(color: Color(0xFFc9d1d9), fontWeight: FontWeight.bold),
        ),
        buttonTheme: ButtonThemeData(
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
          buttonColor: const Color(0xFF238636),
          textTheme: ButtonTextTheme.primary,
        ),
      ),
      home: const PontoScreen(),
    );
  }
}

class PontoScreen extends StatefulWidget {
  const PontoScreen({super.key});

  @override
  State<PontoScreen> createState() => _PontoScreenState();
}

class _PontoScreenState extends State<PontoScreen> {
  List<Map<String, dynamic>> _dailyRecords = [];
  String _totalHoursWorked = '00h 00m';
  String _predictedExit = 'N/A';
  DateTime _selectedDate = DateTime.now();

  @override
  void initState() {
    super.initState();
    // Lógica de inicialização de dados ainda desativada
    // _fetchDataForSelectedDate();
  }

  Future<void> _fetchDataForSelectedDate() async {
    // Lógica desativada por enquanto
  }

  Future<void> _fetchDailyRecords() async {
    // Lógica desativada por enquanto
  }

  void _calculateTotalHours() {
    // Lógica desativada por enquanto
  }

  void _calculatePredictedExit() {
    // Lógica desativada por enquanto
  }

  Future<void> _registrarPonto([DateTime? specificTime]) async {
    // Lógica desativada por enquanto
  }

  Future<void> _updatePonto(String id, TimeOfDay newTime) async {
    // Lógica desativada por enquanto
  }

  Future<void> _deletePonto(String id) async {
    // Lógica desativada por enquanto
  }

  void _goToPreviousDay() {
    setState(() {
      _selectedDate = _selectedDate.subtract(const Duration(days: 1));
    });
    // _fetchDataForSelectedDate();
  }

  void _goToNextDay() {
    setState(() {
      _selectedDate = _selectedDate.add(const Duration(days: 1));
    });
    // _fetchDataForSelectedDate();
  }

  Future<void> _showDateTimePicker() async {
    // Lógica desativada por enquanto
  }

  Future<void> _showEditDialog(String id, DateTime currentTime) async {
    // Lógica desativada por enquanto
  }

  String _formatDate(DateTime date) {
    var formatter = DateFormat("EEEE, d 'de' MMMM", 'pt_BR');
    return formatter.format(date);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Center(
          child: ConstrainedBox(
            constraints: const BoxConstraints(maxWidth: 600),
            child: Container(
              margin: const EdgeInsets.all(16),
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                color: Theme.of(context).colorScheme.surface,
                borderRadius: BorderRadius.circular(24),
                boxShadow: const [
                  BoxShadow(
                    color: Colors.black26,
                    blurRadius: 20,
                    offset: Offset(0, 10),
                  )
                ],
              ),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  _buildNavigation(),
                  const SizedBox(height: 24),
                  _buildInfoCards(),
                  const SizedBox(height: 24),
                  _buildEntryList(),
                  const SizedBox(height: 24),
                  _buildActionButtons(),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildNavigation() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        IconButton(
          icon: const Icon(Icons.chevron_left, size: 28),
          onPressed: _goToPreviousDay,
          tooltip: 'Dia Anterior',
        ),
        Text(
          _formatDate(_selectedDate),
          style: Theme.of(context).textTheme.titleLarge?.copyWith(fontSize: 18),
        ),
        IconButton(
          icon: const Icon(Icons.chevron_right, size: 28),
          onPressed: _goToNextDay,
          tooltip: 'Próximo Dia',
        ),
      ],
    );
  }

  Widget _buildInfoCards() {
    return Row(
      children: [
        Expanded(
          child: _infoCard('Total trabalhado', _totalHoursWorked, const Color(0xFF58a6ff)),
        ),
        const SizedBox(width: 16),
        Expanded(
          child: _infoCard('Saída prevista', _predictedExit, const Color(0xFFa5d6ff)),
        ),
      ],
    );
  }

  Widget _infoCard(String title, String value, Color valueColor) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xFF21262d),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        children: [
          Text(title, style: const TextStyle(color: Color(0xFF8b949e), fontSize: 12)),
          const SizedBox(height: 4),
          Text(value, style: TextStyle(color: valueColor, fontSize: 20, fontWeight: FontWeight.bold)),
        ],
      ),
    );
  }

  Widget _buildEntryList() {
    if (_dailyRecords.isEmpty) {
      return const Padding(
        padding: EdgeInsets.symmetric(vertical: 32.0),
        child: Text(
          'Nenhum ponto registrado para este dia.',
          textAlign: TextAlign.center,
          style: TextStyle(color: Color(0xFF8b949e)),
        ),
      );
    }
    return ListView.builder(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      itemCount: _dailyRecords.length,
      itemBuilder: (context, index) {
        final record = _dailyRecords[index];
        final isCheckIn = record['type'] == 'check-in';
        return _entryTile(
          id: record['id'].toString(),
          time: DateFormat('HH:mm').format(record['timestamp']),
          type: isCheckIn ? 'Entrada' : 'Saída',
          icon: isCheckIn ? Icons.arrow_forward : Icons.arrow_back,
          iconColor: isCheckIn ? Colors.greenAccent : Colors.redAccent,
          onEdit: () => _showEditDialog(record['id'].toString(), record['timestamp']),
          onDelete: () => _deletePonto(record['id'].toString()),
        );
      },
    );
  }

  Widget _entryTile({
    required String id,
    required String time,
    required String type,
    required IconData icon,
    required Color iconColor,
    required VoidCallback onEdit,
    required VoidCallback onDelete,
  }) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      decoration: BoxDecoration(
        color: const Color(0xFF21262d),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          Icon(icon, color: iconColor, size: 24),
          const SizedBox(width: 16),
          Text(time, style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 16)),
          const Spacer(),
          Text(type, style: const TextStyle(color: Color(0xFF8b949e), fontSize: 12)),
          IconButton(
            icon: const Icon(Icons.edit, size: 20),
            onPressed: onEdit,
            tooltip: 'Editar',
          ),
          IconButton(
            icon: const Icon(Icons.delete, size: 20),
            onPressed: onDelete,
            tooltip: 'Deletar',
          ),
        ],
      ),
    );
  }

  Widget _buildActionButtons() {
    return Column(
      children: [
        ElevatedButton(
          onPressed: () => _registrarPonto(),
          style: ElevatedButton.styleFrom(
            backgroundColor: const Color(0xFF238636),
            minimumSize: const Size(double.infinity, 50),
            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
          ),
          child: const Text('Registrar Ponto', style: TextStyle(fontWeight: FontWeight.w600)),
        ),
        const SizedBox(height: 12),
        ElevatedButton(
          onPressed: _showDateTimePicker,
          style: ElevatedButton.styleFrom(
            backgroundColor: const Color(0xFF30363d),
            minimumSize: const Size(double.infinity, 50),
            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
          ),
          child: const Text('Registrar com Data e Hora', style: TextStyle(fontWeight: FontWeight.w600)),
        ),
      ],
    );
  }
}