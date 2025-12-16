import 'package:file_picker/file_picker.dart';
import 'package:client/communication/communication.pbgrpc.dart';
import 'package:client/state/global_state.dart';
import 'package:flutter/material.dart';
import 'package:grpc/grpc_web.dart';
import 'package:provider/provider.dart';

class FileBasedProcessingView extends StatefulWidget {
  const FileBasedProcessingView({super.key});

  @override
  State<FileBasedProcessingView> createState() =>
      _FileBasedProcessingViewState();
}

class _FileBasedProcessingViewState extends State<FileBasedProcessingView> {
  bool _isLoading = true;
  String? _error;
  FileAcquisitionParameters? _params;

  // Form State
  String? _selectedFrameType;
  String? _selectedAcqMode;
  String? _selectedPayload;
  String? _selectedConfigName;
  String? _selectedResultProfile;

  final TextEditingController _remarksController = TextEditingController();

  // File Upload State (Simulated)
  // Map<FrameIdentifier, FileName>
  final Map<String, String> _selectedFiles = {};
  bool _isUploaded = false;

  @override
  void initState() {
    super.initState();
    _fetchParameters();
  }

  @override
  void dispose() {
    _remarksController.dispose();
    super.dispose();
  }

  Future<void> _fetchParameters() async {
    final serverUrl = context.read<GlobalState>().serverUrl;
    final channel = GrpcWebClientChannel.xhr(Uri.parse(serverUrl));
    final client = CommunicationClient(channel);
    try {
      final request = ClientID()..id = context.read<GlobalState>().clientId;
      final response = await client.getFileAcquisitionParameters(request);

      if (mounted) {
        setState(() {
          _params = response;
          _isLoading = false;

          // Set defaults
          if (_params!.frameTypes.isNotEmpty) {
            _selectedFrameType = _params!.frameTypes.first;
          }
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
          // Left Pane: Configuration (60%)
          Expanded(flex: 6, child: _buildInputPanel()),
          const SizedBox(width: 24),
          // Right Pane: Frame Identifiers & File Upload (40%)
          Expanded(flex: 4, child: _buildFileListPanel()),
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
              'Configuration',
              style: Theme.of(context).textTheme.headlineSmall,
            ),
            const SizedBox(height: 8),
            Container(
              padding: const EdgeInsets.all(8.0),
              decoration: BoxDecoration(
                color: Colors.amber[100],
                borderRadius: BorderRadius.circular(4),
                border: Border.all(color: Colors.amber[800]!),
              ),
              child: Row(
                children: [
                  Icon(Icons.warning_amber_rounded, color: Colors.amber[900]),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      'File Based Processing to be Used in Local System Only',
                      style: TextStyle(
                        fontStyle: FontStyle.italic,
                        fontWeight: FontWeight.w600,
                        color: Colors.black87,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 24),
            Expanded(
              child: SingleChildScrollView(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _buildDropdown(
                      label: 'Frame Type',
                      value: _selectedFrameType,
                      items: _params!.frameTypes,
                      onChanged: (val) {
                        setState(() {
                          _selectedFrameType = val;
                          _selectedFiles.clear(); // Clear files on type change
                          _isUploaded = false;
                        });
                      },
                    ),
                    const SizedBox(height: 16),
                    _buildDropdown(
                      label: 'Acquisition Mode',
                      value: _selectedAcqMode,
                      items: _params!.acqModes,
                      onChanged: (val) =>
                          setState(() => _selectedAcqMode = val),
                    ),
                    const SizedBox(height: 16),
                    _buildDropdown(
                      label: 'Payload',
                      value: _selectedPayload,
                      items: _params!.payloads,
                      onChanged: (val) =>
                          setState(() => _selectedPayload = val),
                    ),
                    const SizedBox(height: 16),
                    _buildDropdown(
                      label: 'Config Name',
                      value: _selectedConfigName,
                      items: _params!.configNames,
                      onChanged: (val) =>
                          setState(() => _selectedConfigName = val),
                    ),
                    const SizedBox(height: 16),
                    _buildDropdown(
                      label: 'Result Profile',
                      value: _selectedResultProfile,
                      items: _params!.resultProfiles,
                      onChanged: (val) =>
                          setState(() => _selectedResultProfile = val),
                    ),
                    const SizedBox(height: 24),
                    TextFormField(
                      controller: _remarksController,
                      decoration: const InputDecoration(
                        labelText: 'Remarks',
                        border: OutlineInputBorder(),
                      ),
                      maxLines: 3,
                    ),
                    const SizedBox(height: 32),
                    SizedBox(
                      width: double.infinity,
                      child: FilledButton(
                        onPressed: _isUploaded
                            ? () {
                                // TODO: Implement Start Processing Logic
                                ScaffoldMessenger.of(context).showSnackBar(
                                  const SnackBar(
                                    content: Text('Processing Started'),
                                  ),
                                );
                              }
                            : null,
                        style: FilledButton.styleFrom(
                          padding: const EdgeInsets.symmetric(vertical: 16),
                        ),
                        child: const Text('Start Processing'),
                      ),
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
        return DropdownMenuItem(value: item, child: Text(item));
      }).toList(),
      onChanged: onChanged,
    );
  }

  Widget _buildFileListPanel() {
    // Find Frame Identifiers for current Frame Type
    final frameMap = _params!.frameTypeMap.firstWhere(
      (m) => m.frameType == _selectedFrameType,
      orElse: () => FrameTypeMap(),
    );

    final identifiers = frameMap.frameIdentifiers;

    return Card(
      elevation: 2,
      child: Padding(
        padding: const EdgeInsets.all(24.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Frame Identifiers & Files',
              style: Theme.of(context).textTheme.headlineSmall,
            ),
            const SizedBox(height: 8),
            Text(
              'Selected Frame Type: ${_selectedFrameType ?? "None"}',
              style: Theme.of(
                context,
              ).textTheme.bodyMedium?.copyWith(color: Colors.grey[600]),
            ),
            const SizedBox(height: 24),
            if (identifiers.isEmpty)
              const Center(
                child: Text('No identifiers found for this Frame Type.'),
              )
            else ...[
              Expanded(
                child: ListView.separated(
                  itemCount: identifiers.length,
                  separatorBuilder: (context, index) => const Divider(),
                  itemBuilder: (context, index) {
                    final id = identifiers[index];
                    return _buildFileUploaderRow(id);
                  },
                ),
              ),
              const SizedBox(height: 16),
              SizedBox(
                width: double.infinity,
                child: FilledButton.icon(
                  onPressed: _selectedFiles.isEmpty
                      ? null
                      : () {
                          setState(() {
                            _isUploaded = true;
                          });
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(
                              content: Text('Files Uploaded Successfully'),
                            ),
                          );
                        },
                  icon: const Icon(Icons.upload),
                  label: const Text('Upload All Files'),
                  style: FilledButton.styleFrom(
                    padding: const EdgeInsets.symmetric(vertical: 16),
                  ),
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildFileUploaderRow(String identifier) {
    final fileName = _selectedFiles[identifier] ?? '';

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          // Identifier Name
          SizedBox(
            width: 100,
            child: Text(
              identifier,
              style: const TextStyle(fontWeight: FontWeight.bold),
            ),
          ),
          const SizedBox(width: 16),
          // File Browser / Path Display
          Expanded(
            child: Row(
              children: [
                Expanded(
                  child: Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 12,
                      vertical: 14,
                    ),
                    decoration: BoxDecoration(
                      border: Border.all(color: Colors.grey),
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: Text(
                      fileName.isEmpty ? 'No file selected' : fileName,
                      style: TextStyle(
                        color: fileName.isEmpty ? Colors.grey : Colors.black,
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                  ),
                ),
                const SizedBox(width: 8),
                OutlinedButton.icon(
                  onPressed: () async {
                    FilePickerResult? result = await FilePicker.platform
                        .pickFiles();

                    if (result != null) {
                      setState(() {
                        // Use path if available (Desktop), otherwise name (Web)
                        _selectedFiles[identifier] =
                            result.files.single.path ??
                            result.files.single.name;
                        _isUploaded = false;
                      });
                    }
                  },
                  icon: const Icon(Icons.folder_open),
                  label: const Text('Browse'),
                  style: OutlinedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(
                      vertical: 16,
                      horizontal: 16,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
