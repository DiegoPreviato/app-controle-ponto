import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:app_controle_ponto_frontend/widgets/manual_time_picker_dialog.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Controle de Ponto',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
        useMaterial3: true,
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
  String _message = '';
  List<DateTime> _dailyRecords = [];
  String _totalHoursWorked = '0h 0m 0s'; // New state variable
  DateTime _selectedDate = DateTime.now(); // New state variable for date selection

  @override
  void initState() {
    super.initState();
    _fetchDailyRecords(); // These will now implicitly use _selectedDate
    _fetchTotalHours(); // Call to fetch total hours
  }

  Future<void> _fetchDailyRecords() async {
    setState(() {
      _message = 'Carregando registros do dia...';
    });

    try {
      final dateToFetch = _selectedDate; // Use selected date
      final formattedDate = "${dateToFetch.year}-${dateToFetch.month.toString().padLeft(2, '0')}-${dateToFetch.day.toString().padLeft(2, '0')}";
      final url = Uri.parse('http://127.0.0.1:8080/pontos/$formattedDate');

      print('Fetching records from: $url');

      final response = await http.get(url);

      print('Fetch response status: ${response.statusCode}');
      print('Fetch response body: ${response.body}');

      if (response.statusCode == 200) {
        final List<dynamic> recordsJson = jsonDecode(response.body);
        setState(() {
          _dailyRecords = recordsJson.map((record) => DateTime.parse(record['horario'])).toList();
          _message = 'Registros carregados.';
        });
      } else if (response.statusCode == 204) { // No Content
        setState(() {
          _dailyRecords = [];
          _message = 'Nenhum registro para hoje.';
        });
      }
      else {
        setState(() {
          _message = 'Erro ao carregar registros: ${response.statusCode} ${response.body}';
        });
      }
    } catch (e) {
      setState(() {
        _message = 'Erro de conexão ao carregar registros: $e';
      });
      print('Connection error fetching records: $e');
    }
  }

  // New function to fetch total hours
  Future<void> _fetchTotalHours() async {
    try {
      final dateToFetch = _selectedDate; // Use selected date
      final formattedDate = "${dateToFetch.year}-${dateToFetch.month.toString().padLeft(2, '0')}-${dateToFetch.day.toString().padLeft(2, '0')}";
      final url = Uri.parse('http://127.0.0.1:8080/pontos/$formattedDate/total-horas');

      print('Fetching total hours from: $url');

      final response = await http.get(url);

      print('Total hours response status: ${response.statusCode}');
      print('Total hours response body: ${response.body}');

      if (response.statusCode == 200) {
        final Map<String, dynamic> data = jsonDecode(response.body);
        setState(() {
          _totalHoursWorked = data['total_trabalhado'] ?? '0h 0m 0s';
        });
      } else {
        print('Erro ao carregar total de horas: ${response.statusCode} ${response.body}');
        setState(() {
          _totalHoursWorked = 'Erro';
        });
      }
    } catch (e) {
      print('Connection error fetching total hours: $e');
      setState(() {
        _totalHoursWorked = 'Erro';
      });
    }
  }

  Future<void> _registrarPonto([DateTime? specificTime]) async {
    setState(() {
      _message = 'Registrando ponto...';
    });

    http.Response? response; // Declare response here, initialized to null

    try {
      final now = specificTime ?? DateTime.now();
      final formattedDate = "${now.year}-${now.month.toString().padLeft(2, '0')}-${now.day.toString().padLeft(2, '0')}";
      final hour = now.hour;
      final minute = now.minute;

      final url = Uri.parse('http://127.0.0.1:8080/registrar-ponto');
      final headers = <String, String>{
        'Content-Type': 'application/json; charset=UTF-8',
      };
      final body = jsonEncode(<String, dynamic>{
        'data': formattedDate,
        'hora': hour,
        'minuto': minute,
      });

      print('Sending request to: $url');
      print('Headers: $headers');
      print('Body: $body');

      response = await http.post( // Assign to the already declared response
        url,
        headers: headers,
        body: body,
      );

      print('Response status: ${response?.statusCode}'); // Null-safe access
      print('Response body: ${response?.body}'); // Null-safe access

      if (response != null && response.statusCode == 201) { // Null check before accessing statusCode
        setState(() {
          _message = 'Ponto registrado com sucesso!';
        });
        _fetchDailyRecords(); // Refresh records after successful registration
        _fetchTotalHours(); // Refresh total hours after successful registration
      } else {
        setState(() {
          _message = 'Erro ao registrar ponto: ${response?.statusCode} ${response?.body}'; // Null-safe access
        });
      }
    } catch (e) {
      setState(() {
        _message = 'Erro de conexão: $e';
      });
      print('Connection error: $e');
    }
  }

  void _goToPreviousDay() {
    setState(() {
      _selectedDate = _selectedDate.subtract(const Duration(days: 1));
    });
    _fetchDailyRecords();
    _fetchTotalHours();
  }

  void _goToNextDay() {
    setState(() {
      _selectedDate = _selectedDate.add(const Duration(days: 1));
    });
    _fetchDailyRecords();
    _fetchTotalHours();
  }

  Future<void> _selectDate(BuildContext context) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: _selectedDate,
      firstDate: DateTime(2000), // Arbitrary start date
      lastDate: DateTime(2101),  // Arbitrary end date
    );
    if (picked != null && picked != _selectedDate) {
      setState(() {
        _selectedDate = picked;
      });
      _fetchDailyRecords();
      _fetchTotalHours();
    }
  }

  Future<void> _registerPontoAtSpecificTime(BuildContext context) async {
    final DateTime? selectedDateTime = await showDialog<DateTime>(
      context: context,
      builder: (BuildContext context) {
        return ManualTimePickerDialog(initialDate: DateTime.now());
      },
    );

    if (selectedDateTime != null) {
      await _registrarPonto(selectedDateTime);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        title: const Text('Controle de Ponto'),
        actions: [ // Add actions for total hours
          Padding(
            padding: const EdgeInsets.only(right: 16.0),
            child: Center(
              child: Text(
                'Total: $_totalHoursWorked',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                  color: Theme.of(context).colorScheme.onPrimaryContainer, // Adjust color as needed
                ),
              ),
            ),
          ),
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () => _registerPontoAtSpecificTime(context),
          ),
        ],
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.start, // Align to start for records
          crossAxisAlignment: CrossAxisAlignment.stretch, // Stretch to fill width
          children: <Widget>[
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                IconButton(
                  icon: const Icon(Icons.arrow_left),
                  onPressed: _goToPreviousDay,
                ),
                GestureDetector( // New: Wrap Text with GestureDetector
                  onTap: () => _selectDate(context), // Call _selectDate on tap
                  child: Text(
                    '${_selectedDate.day.toString().padLeft(2, '0')}/${_selectedDate.month.toString().padLeft(2, '0')}/${_selectedDate.year}',
                    style: Theme.of(context).textTheme.headlineMedium,
                  ),
                ),
                IconButton(
                  icon: const Icon(Icons.arrow_right),
                  onPressed: _goToNextDay,
                ),
              ],
            ),
            const SizedBox(height: 20), // Add some spacing
            ElevatedButton(
              onPressed: _registrarPonto,
              child: const Text('Registrar Ponto'),
            ),
            const SizedBox(height: 10),
            ElevatedButton(
              onPressed: () => _registerPontoAtSpecificTime(context),
              child: const Text('Registrar Ponto (Data Específica)'),
            ),
            const SizedBox(height: 20),
            Text(
              _message,
              style: Theme.of(context).textTheme.headlineSmall,
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 20),
            Text(
              'Registros do Dia:',
              style: Theme.of(context).textTheme.headlineSmall,
            ),
            const SizedBox(height: 10),
            Expanded( // Use Expanded to allow ListView to take available space
              child: _dailyRecords.isEmpty
                  ? Center(child: Text('Nenhum registro para exibir.'))
                  : ListView.builder(
                      itemCount: _dailyRecords.length,
                      itemBuilder: (context, index) {
                        final record = _dailyRecords[index];
                        final type = index % 2 == 0 ? 'Entrada' : 'Saída';
                        final localRecord = record.toLocal();
                        final formattedTime = '${localRecord.hour.toString().padLeft(2, '0')}:${localRecord.minute.toString().padLeft(2, '0')}';
                        return Padding(
                          padding: const EdgeInsets.symmetric(vertical: 4.0),
                          child: Text('$type: $formattedTime', style: Theme.of(context).textTheme.bodyLarge),
                        );
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }
}
