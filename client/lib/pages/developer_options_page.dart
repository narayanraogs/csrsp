import 'package:client/communication/communication.pbgrpc.dart';
import 'package:client/state/global_state.dart';
import 'package:client/state/log_state.dart';
import 'package:flutter/material.dart';
import 'package:grpc/grpc_web.dart';
import 'package:provider/provider.dart';

class DeveloperOptionsView extends StatefulWidget {
  const DeveloperOptionsView({super.key});

  @override
  State<DeveloperOptionsView> createState() => _DeveloperOptionsViewState();
}

class _DeveloperOptionsViewState extends State<DeveloperOptionsView> {
  bool _isLoading = false;
  DeveloperOptions? _options;

  // Controllers
  bool _autoArchival = false;
  String _logLevel = 'INFO';
  bool _parallelProcessing = false;
  String _encryptionMode = 'Internal Index';
  final TextEditingController _endProcessIdController = TextEditingController();
  final TextEditingController _maxThreadsController = TextEditingController();
  final TextEditingController _bufferLengthController = TextEditingController();

  final List<String> _logLevels = ['DEBUG', 'INFO', 'WARN', 'ERROR'];
  final List<String> _encryptionModes = ['Internal Index', 'External Key'];

  @override
  void initState() {
    super.initState();
    _fetchOptions();
  }

  @override
  void dispose() {
    _endProcessIdController.dispose();
    _maxThreadsController.dispose();
    _bufferLengthController.dispose();
    super.dispose();
  }

  Future<void> _fetchOptions() async {
    setState(() => _isLoading = true);

    try {
      final globalState = context.read<GlobalState>();
      final channel = GrpcWebClientChannel.xhr(
        Uri.parse(globalState.serverUrl),
      );
      final client = CommunicationClient(channel);
      final request = ClientID()..id = globalState.clientId;

      final response = await client.getDeveloperOptions(request);

      if (mounted) {
        setState(() {
          _options = response;
          _autoArchival = response.autoArchival;
          _logLevel = response.logLevel;
          _parallelProcessing = response.enableParallelProcessing;

          if (_encryptionModes.contains(response.encryptionMode)) {
            _encryptionMode = response.encryptionMode;
          } else {
            _encryptionMode = _encryptionModes.first;
          }

          _endProcessIdController.text = response.endProcessID.toString();
          _maxThreadsController.text = response.maxThreads.toString();
          _bufferLengthController.text = response.bufferLength.toString();
          _isLoading = false;
        });
      }
    } catch (e) {
      debugPrint('Error fetching developer options: $e');
      if (mounted) {
        setState(() => _isLoading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Failed to load developer options: $e'),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  Future<void> _saveOptions() async {
    setState(() => _isLoading = true);

    try {
      final globalState = context.read<GlobalState>();
      final channel = GrpcWebClientChannel.xhr(
        Uri.parse(globalState.serverUrl),
      );
      final client = CommunicationClient(channel);

      // Parse numeric values
      final endProcessID = int.tryParse(_endProcessIdController.text) ?? -1;
      final maxThreads = int.tryParse(_maxThreadsController.text) ?? 1;
      final bufferLength = int.tryParse(_bufferLengthController.text) ?? 1024;

      final request = DeveloperOptions()
        ..autoArchival = _autoArchival
        ..logLevel = _logLevel
        ..enableParallelProcessing = _parallelProcessing
        ..encryptionMode = _encryptionMode
        ..endProcessID = endProcessID
        ..maxThreads = maxThreads
        ..bufferLength = bufferLength;

      final response = await client.setDeveloperOptions(request);

      if (mounted) {
        context.read<LogState>().addLog(
          response.ok
              ? 'Developer options saved successfully'
              : 'Failed to save: ${response.message}',
          type: response.ok ? LogType.success : LogType.error,
          showSnackbar: true,
          context: context,
        );
      }
    } catch (e) {
      debugPrint('Error saving developer options: $e');
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e'), backgroundColor: Colors.red),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading && _options == null) {
      return const Center(child: CircularProgressIndicator());
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text('Developer Options'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: _isLoading ? null : _fetchOptions,
          ),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildSectionHeader('General Settings'),
            Card(
              elevation: 2,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
              child: Column(
                children: [
                  SwitchListTile(
                    title: const Text('Automatic Archival'),
                    subtitle: const Text(
                      'Enable automatic data archival after acquisition',
                    ),
                    value: _autoArchival,
                    onChanged: (val) => setState(() => _autoArchival = val),
                  ),
                  const Divider(height: 1),
                  SwitchListTile(
                    title: const Text('Parallel Processing'),
                    subtitle: const Text('Enable concurrent data processing'),
                    value: _parallelProcessing,
                    onChanged: (val) =>
                        setState(() => _parallelProcessing = val),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 24),

            _buildSectionHeader('Logging & Security'),
            Card(
              elevation: 2,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
              child: Padding(
                padding: const EdgeInsets.all(16.0),
                child: Column(
                  children: [
                    DropdownButtonFormField<String>(
                      value: _logLevels.contains(_logLevel)
                          ? _logLevel
                          : _logLevels.first,
                      decoration: const InputDecoration(
                        labelText: 'Log Level',
                        border: OutlineInputBorder(),
                        prefixIcon: Icon(Icons.analytics),
                      ),
                      items: _logLevels.map((level) {
                        return DropdownMenuItem(
                          value: level,
                          child: Text(level),
                        );
                      }).toList(),
                      onChanged: (val) {
                        if (val != null) setState(() => _logLevel = val);
                      },
                    ),
                    const SizedBox(height: 16),
                    DropdownButtonFormField<String>(
                      value: _encryptionMode,
                      decoration: const InputDecoration(
                        labelText: 'Encryption Mode',
                        border: OutlineInputBorder(),
                        prefixIcon: Icon(Icons.security),
                      ),
                      items: _encryptionModes.map((mode) {
                        return DropdownMenuItem(value: mode, child: Text(mode));
                      }).toList(),
                      onChanged: (val) {
                        if (val != null) setState(() => _encryptionMode = val);
                      },
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 24),

            _buildSectionHeader('Performance & Limits'),
            Card(
              elevation: 2,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
              child: Padding(
                padding: const EdgeInsets.all(16.0),
                child: Column(
                  children: [
                    TextField(
                      controller: _endProcessIdController,
                      decoration: const InputDecoration(
                        labelText: 'End Process ID',
                        helperText:
                            'ID to signal process termination (-1 to disable)',
                        border: OutlineInputBorder(),
                        prefixIcon: Icon(Icons.stop_circle),
                      ),
                      keyboardType: TextInputType.number,
                    ),
                    const SizedBox(height: 16),
                    TextField(
                      controller: _maxThreadsController,
                      decoration: const InputDecoration(
                        labelText: 'Max Threads',
                        helperText: 'Maximum number of concurrent threads',
                        border: OutlineInputBorder(),
                        prefixIcon: Icon(Icons.speed),
                      ),
                      keyboardType: TextInputType.number,
                    ),
                    const SizedBox(height: 16),
                    TextField(
                      controller: _bufferLengthController,
                      decoration: const InputDecoration(
                        labelText: 'Buffer Length',
                        helperText: 'Size of processing buffer in bytes',
                        border: OutlineInputBorder(),
                        prefixIcon: Icon(Icons.memory),
                      ),
                      keyboardType: TextInputType.number,
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 32),

            SizedBox(
              width: double.infinity,
              height: 50,
              child: ElevatedButton.icon(
                onPressed: _isLoading ? null : _saveOptions,
                icon: _isLoading
                    ? const SizedBox(
                        width: 20,
                        height: 20,
                        child: CircularProgressIndicator(
                          color: Colors.white,
                          strokeWidth: 2,
                        ),
                      )
                    : const Icon(Icons.save),
                label: const Text('Save Configuration'),
                style: ElevatedButton.styleFrom(
                  elevation: 4,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(8),
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSectionHeader(String title) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12.0, left: 4.0),
      child: Text(
        title,
        style: TextStyle(
          fontSize: 18,
          fontWeight: FontWeight.bold,
          color: Theme.of(context).primaryColor,
        ),
      ),
    );
  }
}
