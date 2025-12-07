import 'package:client/communication/communication.pbgrpc.dart';
import 'package:client/state/global_state.dart';
import 'package:flutter/material.dart';
import 'package:grpc/grpc_web.dart';
import 'package:provider/provider.dart';

class AcquireDataView extends StatefulWidget {
  const AcquireDataView({super.key});

  @override
  State<AcquireDataView> createState() => _AcquireDataViewState();
}

class _AcquireDataViewState extends State<AcquireDataView> {
  bool _isLoading = true;
  String? _error;
  AcquisitionParameters? _params;

  // Form State
  String? _selectedAcqMode;
  String? _selectedPayload;
  String? _selectedConfigName;
  String? _selectedResultProfile;
  String _acqType = 'Frame'; // Default to Frame

  final TextEditingController _frameCountController = TextEditingController(text:"50000");
  final TextEditingController _timeController = TextEditingController(text:"30");
  final TextEditingController _remarksController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _fetchParameters();
  }

  @override
  void dispose() {
    _frameCountController.dispose();
    _timeController.dispose();
    _remarksController.dispose();
    super.dispose();
  }

  Future<void> _fetchParameters() async {
    final serverUrl = context.read<GlobalState>().serverUrl;
    final channel = GrpcWebClientChannel.xhr(Uri.parse(serverUrl));
    final client = CommunicationClient(channel);

    try {
      final request = ClientID()..id = context.read<GlobalState>().clientId;
      final response = await client.getAcquisitionParameters(request);

      if (mounted) {
        setState(() {
          _params = response;
          _isLoading = false;

          // Set defaults if available
          if (_params!.acqModes.isNotEmpty) {
            _selectedAcqMode = _params!.acqModes.first;
          }
          if (_params!.payloads.isNotEmpty) {
            _selectedPayload = _params!.payloads.first;
          }
          if (_params!.configNames.isNotEmpty) {
            _selectedConfigName = _params!.configNames.first;
          }
          if (_params!.resultProfiles.isNotEmpty) {
            _selectedResultProfile = _params!.resultProfiles.first;
          }
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _error = e.toString();
          _isLoading = false;
        });
      }
    } finally {
      channel.shutdown();
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_error != null) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.error, color: Colors.red, size: 48),
            const SizedBox(height: 16),
            Text('Error loading parameters: $_error'),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: () {
                setState(() {
                  _isLoading = true;
                  _error = null;
                });
                _fetchParameters();
              },
              child: const Text('Retry'),
            ),
          ],
        ),
      );
    }

    if (_params == null) {
      return const Center(child: Text('No parameters received.'));
    }

    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Portion 1: Inputs
          Expanded(
            flex: 2,
            child: _buildInputPanel(),
          ),
          const SizedBox(width: 24),
          // Portion 2: DAS Map
          Expanded(
            flex: 1,
            child: _buildDasPanel(),
          ),
        ],
      ),
    );
  }

  Widget _buildInputPanel() {
    return Card(
      elevation: 2,
      child: Padding(
        padding: const EdgeInsets.all(24.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Acquisition Configuration',
              style: Theme.of(context).textTheme.headlineSmall,
            ),
            const SizedBox(height: 24),
            Expanded(
              child: SingleChildScrollView(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Dropdowns
                    _buildDropdown(
                      label: 'Acquisition Mode',
                      value: _selectedAcqMode,
                      items: _params!.acqModes,
                      onChanged: (val) => setState(() => _selectedAcqMode = val),
                    ),
                    const SizedBox(height: 16),
                    _buildDropdown(
                      label: 'Payload',
                      value: _selectedPayload,
                      items: _params!.payloads,
                      onChanged: (val) => setState(() => _selectedPayload = val),
                    ),
                    const SizedBox(height: 16),
                    _buildDropdown(
                      label: 'Config Name',
                      value: _selectedConfigName,
                      items: _params!.configNames,
                      onChanged: (val) => setState(() => _selectedConfigName = val),
                    ),
                    const SizedBox(height: 16),
                    _buildDropdown(
                      label: 'Result Profile',
                      value: _selectedResultProfile,
                      items: _params!.resultProfiles,
                      onChanged: (val) => setState(() => _selectedResultProfile = val),
                    ),
                    const SizedBox(height: 24),

                    // Acquisition Type
                    Text('Acquisition Type',
                        style: Theme.of(context).textTheme.titleMedium),
                    const SizedBox(height: 8),
                    Column(
                      children: [
                        RadioListTile<String>(
                          title: const Text('Frame'),
                          value: 'Frame',
                          groupValue: _acqType,
                          onChanged: (val) => setState(() => _acqType = val!),
                        ),
                        if (_acqType == 'Frame')
                          Padding(
                            padding: const EdgeInsets.only(left: 32.0, bottom: 8.0),
                            child: TextFormField(
                              controller: _frameCountController,
                              decoration: const InputDecoration(
                                labelText: 'Number of Frames',
                                border: OutlineInputBorder(),
                                isDense: true,
                              ),
                              keyboardType: TextInputType.number,
                            ),
                          ),
                        RadioListTile<String>(
                          title: const Text('Time'),
                          value: 'Time',
                          groupValue: _acqType,
                          onChanged: (val) => setState(() => _acqType = val!),
                        ),
                        if (_acqType == 'Time')
                          Padding(
                            padding: const EdgeInsets.only(left: 32.0, bottom: 8.0),
                            child: TextFormField(
                              controller: _timeController,
                              decoration: const InputDecoration(
                                labelText: 'Time (seconds)',
                                border: OutlineInputBorder(),
                                isDense: true,
                              ),
                              keyboardType: TextInputType.number,
                            ),
                          ),
                        RadioListTile<String>(
                          title: const Text('User Defined'),
                          value: 'UserDefined',
                          groupValue: _acqType,
                          onChanged: (val) => setState(() => _acqType = val!),
                        ),
                      ],
                    ),
                    const SizedBox(height: 24),

                    // Remarks
                    TextFormField(
                      controller: _remarksController,
                      decoration: const InputDecoration(
                        labelText: 'Remarks',
                        border: OutlineInputBorder(),
                      ),
                      maxLines: 3,
                    ),
                    const SizedBox(height: 32),

                    // Buttons
                    Row(
                      children: [
                        Expanded(
                          child: OutlinedButton(
                            onPressed: () {
                              // TODO: Implement Configure Logic
                              ScaffoldMessenger.of(context).showSnackBar(
                                  const SnackBar(content: Text('System Configured')));
                            },
                            style: OutlinedButton.styleFrom(
                              padding: const EdgeInsets.symmetric(vertical: 16),
                            ),
                            child: const Text('Configure System'),
                          ),
                        ),
                        const SizedBox(width: 16),
                        Expanded(
                          child: FilledButton(
                            onPressed: () {
                              // TODO: Implement Start Logic
                              ScaffoldMessenger.of(context).showSnackBar(
                                  const SnackBar(content: Text('Acquisition Started')));
                            },
                            style: FilledButton.styleFrom(
                              padding: const EdgeInsets.symmetric(vertical: 16),
                            ),
                            child: const Text('Start Acquisition'),
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildDropdown({
    required String label,
    required String? value,
    required List<String> items,
    required ValueChanged<String?> onChanged,
  }) {
    return DropdownButtonFormField<String>(
      value: value,
      decoration: InputDecoration(
        labelText: label,
        border: const OutlineInputBorder(),
      ),
      items: items.map((item) {
        return DropdownMenuItem(
          value: item,
          child: Text(item),
        );
      }).toList(),
      onChanged: onChanged,
    );
  }

  Widget _buildDasPanel() {
    // Find the DAS map for the selected mode
    final dasMapEntry = _params!.dasMap.firstWhere(
      (entry) => entry.acqMode == _selectedAcqMode,
      orElse: () => AcqDASMap(), // Return empty if not found
    );

    return Card(
      elevation: 2,
      child: Padding(
        padding: const EdgeInsets.all(24.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Hardware Status',
              style: Theme.of(context).textTheme.headlineSmall,
            ),
            const SizedBox(height: 8),
            Text(
              'Mode: ${_selectedAcqMode ?? "None"}',
              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                    color: Colors.grey[600],
                  ),
            ),
            const SizedBox(height: 24),
            if (dasMapEntry.dasDetails.isEmpty)
              const Center(
                child: Text('No hardware configuration found for this mode.'),
              )
            else
              Expanded(
                child: ListView.builder(
                  itemCount: dasMapEntry.dasDetails.length,
                  itemBuilder: (context, index) {
                    final detail = dasMapEntry.dasDetails[index];
                    return _buildDasCard(detail);
                  },
                ),
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildDasCard(AcqDasDetails detail) {
    // Status is currently mocked as "Not Connected" (Red)
    const statusColor = Colors.red;
    const statusText = "Not Connected";

    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      shape: RoundedRectangleBorder(
        side: const BorderSide(color: statusColor, width: 2),
        borderRadius: BorderRadius.circular(8),
      ),
      child: ListTile(
        leading: const Icon(Icons.dns, size: 32, color: statusColor),
        title: Text(
          detail.dasName,
          style: const TextStyle(fontWeight: FontWeight.bold),
        ),
        subtitle: Text('DPU: ${detail.dpuNumber}'),
        trailing: Chip(
          label: const Text(
            statusText,
            style: TextStyle(color: Colors.white),
          ),
          backgroundColor: statusColor,
        ),
      ),
    );
  }
}