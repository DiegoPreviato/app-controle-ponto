import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class ManualTimePickerDialog extends StatefulWidget {
  final DateTime initialDate;

  const ManualTimePickerDialog({Key? key, required this.initialDate}) : super(key: key);

  @override
  State<ManualTimePickerDialog> createState() => _ManualTimePickerDialogState();
}

class _ManualTimePickerDialogState extends State<ManualTimePickerDialog> {
  late DateTime _selectedDate;
  final TextEditingController _hourController = TextEditingController();
  final TextEditingController _minuteController = TextEditingController();
  final _formKey = GlobalKey<FormState>();

  @override
  void initState() {
    super.initState();
    _selectedDate = widget.initialDate;
    _hourController.text = _selectedDate.hour.toString().padLeft(2, '0');
    _minuteController.text = _selectedDate.minute.toString().padLeft(2, '0');
  }

  @override
  void dispose() {
    _hourController.dispose();
    _minuteController.dispose();
    super.dispose();
  }

  Future<void> _pickDate() async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: _selectedDate,
      firstDate: DateTime(2000),
      lastDate: DateTime(2101),
    );
    if (picked != null && picked != _selectedDate) {
      setState(() {
        _selectedDate = DateTime(
          picked.year,
          picked.month,
          picked.day,
          _selectedDate.hour,
          _selectedDate.minute,
        );
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Selecionar Data e Hora'),
      content: Form(
        key: _formKey,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            ListTile(
              title: Text(
                'Data: ${_selectedDate.day.toString().padLeft(2, '0')}/${_selectedDate.month.toString().padLeft(2, '0')}/${_selectedDate.year}',
              ),
              trailing: const Icon(Icons.calendar_today),
              onTap: _pickDate,
            ),
            Row(
              children: [
                Expanded(
                  child: TextFormField(
                    controller: _hourController,
                    keyboardType: TextInputType.number,
                    inputFormatters: [
                      FilteringTextInputFormatter.digitsOnly,
                      LengthLimitingTextInputFormatter(2),
                    ],
                    decoration: const InputDecoration(
                      labelText: 'Hora',
                      hintText: 'HH',
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Campo obrigatório';
                      }
                      final hour = int.tryParse(value);
                      if (hour == null || hour < 0 || hour > 23) {
                        return 'Hora inválida (0-23)';
                      }
                      return null;
                    },
                  ),
                ),
                const Text(' : '),
                Expanded(
                  child: TextFormField(
                    controller: _minuteController,
                    keyboardType: TextInputType.number,
                    inputFormatters: [
                      FilteringTextInputFormatter.digitsOnly,
                      LengthLimitingTextInputFormatter(2),
                    ],
                    decoration: const InputDecoration(
                      labelText: 'Minuto',
                      hintText: 'MM',
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Campo obrigatório';
                      }
                      final minute = int.tryParse(value);
                      if (minute == null || minute < 0 || minute > 59) {
                        return 'Minuto inválido (0-59)';
                      }
                      return null;
                    },
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: const Text('Cancelar'),
        ),
        TextButton(
          onPressed: () {
            if (_formKey.currentState!.validate()) {
              final hour = int.parse(_hourController.text);
              final minute = int.parse(_minuteController.text);
              final resultDateTime = DateTime(
                _selectedDate.year,
                _selectedDate.month,
                _selectedDate.day,
                hour,
                minute,
              );
              Navigator.of(context).pop(resultDateTime);
            }
          },
          child: const Text('OK'),
        ),
      ],
    );
  }
}
