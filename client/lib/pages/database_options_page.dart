import 'package:client/communication/communication.pbgrpc.dart';
import 'package:client/state/global_state.dart';
import 'package:client/state/log_state.dart';
import 'package:flutter/material.dart';
import 'package:grpc/grpc_web.dart';
import 'package:provider/provider.dart';

class DatabaseOptionsView extends StatefulWidget {
  const DatabaseOptionsView({super.key});

  @override
  State<DatabaseOptionsView> createState() => _DatabaseOptionsViewState();
}

class _DatabaseOptionsViewState extends State<DatabaseOptionsView> {
  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 3,
      child: Scaffold(
        appBar: AppBar(
          title: const Text('Database Options'),
          bottom: const TabBar(
            tabs: [
              Tab(icon: Icon(Icons.biotech), text: 'Test Phase'),
              Tab(icon: Icon(Icons.settings_ethernet), text: 'DAS IP'),
              Tab(icon: Icon(Icons.comment), text: 'Remarks'),
            ],
            indicatorColor: Colors.white,
            labelColor: Colors.white,
            unselectedLabelColor: Colors.white70,
          ),
        ),
        body: const TabBarView(
          children: [_TestPhaseTab(), _DasIpTab(), _RemarksTab()],
        ),
      ),
    );
  }
}

// -----------------------------------------------------------------------------
// 1. Test Phase Tab
// -----------------------------------------------------------------------------
class _TestPhaseTab extends StatefulWidget {
  const _TestPhaseTab();

  @override
  State<_TestPhaseTab> createState() => _TestPhaseTabState();
}

class _TestPhaseTabState extends State<_TestPhaseTab> {
  List<String> _availablePhases = [];
  String? _selectedPhase;
  final TextEditingController _newPhaseController = TextEditingController();
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _fetchTestPhases();
  }

  @override
  void dispose() {
    _newPhaseController.dispose();
    super.dispose();
  }

  Future<void> _fetchTestPhases() async {
    setState(() {
      _isLoading = true;
    });

    try {
      final globalState = context.read<GlobalState>();
      final channel = GrpcWebClientChannel.xhr(
        Uri.parse(globalState.serverUrl),
      );
      final client = CommunicationClient(channel);
      final request = ClientID()..id = globalState.clientId;

      // We use getAllTestPhases as per proto definition
      final response = await client.getAllTestPhases(request);

      if (mounted) {
        setState(() {
          _availablePhases = response.testPhases;
          // Ideally we would fetch current phase from server too.
          // For now, if we have phases and none selected, pick first.
          if (_selectedPhase == null && _availablePhases.isNotEmpty) {
            _selectedPhase = _availablePhases.first;
          }
          _isLoading = false;
        });
      }
    } catch (e) {
      debugPrint('Error fetching test phases: $e');
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Failed to load test phases: $e')),
        );
      }
    }
  }

  Future<void> _addNewPhase() async {
    final newPhase = _newPhaseController.text.trim();
    if (newPhase.isEmpty) return;

    if (_availablePhases.contains(newPhase)) {
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(const SnackBar(content: Text('Phase already exists.')));
      return;
    }

    setState(() {
      _isLoading = true;
    });

    try {
      final globalState = context.read<GlobalState>();
      final channel = GrpcWebClientChannel.xhr(
        Uri.parse(globalState.serverUrl),
      );
      final client = CommunicationClient(channel);

      final request = TestPhaseRequest()
        ..id = globalState.clientId
        ..testPhase = newPhase;

      final response = await client.addTestPhase(request);

      if (mounted) {
        if (response.ok) {
          _newPhaseController.clear();
          await _fetchTestPhases(); // Refresh list
          await _fetchTestPhases(); // Refresh list
          context.read<LogState>().addLog(
            response.message,
            type: LogType.success,
            showSnackbar: true,
            context: context,
          );
        } else {
          context.read<LogState>().addLog(
            'Error: ${response.message}',
            type: LogType.error,
            showSnackbar: true,
            context: context,
          );
        }
      }
    } catch (e) {
      debugPrint('Error adding test phase: $e');
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('Failed to add test phase: $e')));
      }
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  Future<void> _setActivePhase() async {
    if (_selectedPhase == null) return;

    setState(() {
      _isLoading = true;
    });

    try {
      final globalState = context.read<GlobalState>();
      final channel = GrpcWebClientChannel.xhr(
        Uri.parse(globalState.serverUrl),
      );
      final client = CommunicationClient(channel);

      final request = TestPhaseRequest()
        ..id = globalState.clientId
        ..testPhase = _selectedPhase!;

      final response = await client.selectTestPhase(request);

      if (mounted) {
        context.read<LogState>().addLog(
          response.ok
              ? response.message
              : 'Failed to update phase: ${response.message}',
          type: response.ok ? LogType.success : LogType.error,
          showSnackbar: true,
          context: context,
        );
      }
    } catch (e) {
      debugPrint('Error setting active test phase: $e');
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e'), backgroundColor: Colors.red),
        );
      }
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading && _availablePhases.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    return SingleChildScrollView(
      padding: const EdgeInsets.all(24.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          _buildCard(
            title: 'Select Existing Test Phase',
            child: Column(
              children: [
                DropdownButtonFormField<String>(
                  value: _selectedPhase,
                  decoration: const InputDecoration(
                    labelText: 'Exisiting Phase',
                    border: OutlineInputBorder(),
                    prefixIcon: Icon(Icons.science),
                  ),
                  items: _availablePhases.map((phase) {
                    return DropdownMenuItem(value: phase, child: Text(phase));
                  }).toList(),
                  onChanged: (val) {
                    setState(() {
                      _selectedPhase = val;
                    });
                  },
                ),
                const SizedBox(height: 16),
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton.icon(
                    onPressed: _isLoading ? null : _setActivePhase,
                    icon: const Icon(Icons.check_circle),
                    label: const Text('Update Active Phase'),
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 24),
          _buildCard(
            title: 'Create New Test Phase',
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _newPhaseController,
                    decoration: const InputDecoration(
                      labelText: 'New Phase Name',
                      hintText: 'e.g. Acoustic Testing',
                      border: OutlineInputBorder(),
                    ),
                  ),
                ),
                const SizedBox(width: 16),
                ElevatedButton(
                  onPressed: _isLoading ? null : _addNewPhase,
                  style: ElevatedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(
                      vertical: 16,
                      horizontal: 24,
                    ),
                  ),
                  child: _isLoading
                      ? const SizedBox(
                          width: 20,
                          height: 20,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            color: Colors.white,
                          ),
                        )
                      : const Text('Create'),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCard({required String title, required Widget child}) {
    return Card(
      elevation: 4,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              title,
              style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            ),
            const Divider(),
            const SizedBox(height: 16),
            child,
          ],
        ),
      ),
    );
  }
}

// -----------------------------------------------------------------------------
// 2. DAS IP Tab
// -----------------------------------------------------------------------------
class _DasIpTab extends StatefulWidget {
  const _DasIpTab();

  @override
  State<_DasIpTab> createState() => _DasIpTabState();
}

class _DasIpTabState extends State<_DasIpTab> {
  List<DASIPAddress> _dasList = [];
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _fetchDasIps();
  }

  Future<void> _fetchDasIps() async {
    setState(() {
      _isLoading = true;
    });

    try {
      final globalState = context.read<GlobalState>();
      final channel = GrpcWebClientChannel.xhr(
        Uri.parse(globalState.serverUrl),
      );
      final client = CommunicationClient(channel);
      final request = ClientID()..id = globalState.clientId;

      final response = await client.getDASIPAddresses(request);

      if (mounted) {
        setState(() {
          _dasList = response.dasIPAddresses;
          _isLoading = false;
        });
      }
    } catch (e) {
      debugPrint('Error fetching DAS IPs: $e');
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('Failed to load DAS IPs: $e')));
      }
    }
  }

  Future<void> _saveDasIp(DASIPAddress das, String newIp) async {
    if (newIp.isEmpty) return;

    try {
      final globalState = context.read<GlobalState>();
      final channel = GrpcWebClientChannel.xhr(
        Uri.parse(globalState.serverUrl),
      );
      final client = CommunicationClient(channel);

      // Create a new object for the request
      final request = DASIPAddress()
        ..name = das.name
        ..ipAddress = newIp;

      final response = await client.changeDASIPAddress(request);

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              response.ok
                  ? response.message
                  : 'Failed to update IP: ${response.message}',
            ),
            backgroundColor: response.ok ? Colors.green : Colors.red,
          ),
        );
        if (response.ok) {
          // Update local state to reflect the change firmly
          // (Though the text field should already show it, this syncs our local model)
          setState(() {
            das.ipAddress = newIp;
          });
        }
      }
    } catch (e) {
      debugPrint('Error changing DAS IP: $e');
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e'), backgroundColor: Colors.red),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading && _dasList.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_dasList.isEmpty) {
      return const Center(child: Text("No DAS units found."));
    }

    return ListView.separated(
      padding: const EdgeInsets.all(24.0),
      itemCount: _dasList.length,
      separatorBuilder: (_, __) => const SizedBox(height: 16),
      itemBuilder: (context, index) {
        final das = _dasList[index];
        // We initialize the controller with the current IP.
        // Note: Creating a controller inside build is generally not best practice for performance
        // if the list acts weird (e.g. losing focus), but for this simple list it often suffices.
        // A better approach is maintaining a map of controllers, but let's stick to simple first
        // as we don't expect frequent full-rebuilds while typing.
        // To be safe against losing focus on setState, we can use a key or just be careful.
        // Let's use a Key for the Card to help frameworks.
        return _DasIpCard(
          key: ValueKey(das.name),
          das: das,
          onSave: (newIp) => _saveDasIp(das, newIp),
        );
      },
    );
  }
}

class _DasIpCard extends StatefulWidget {
  final DASIPAddress das;
  final Function(String) onSave;

  const _DasIpCard({required Key key, required this.das, required this.onSave})
    : super(key: key);

  @override
  State<_DasIpCard> createState() => _DasIpCardState();
}

class _DasIpCardState extends State<_DasIpCard> {
  late TextEditingController _controller;

  @override
  void initState() {
    super.initState();
    _controller = TextEditingController(text: widget.das.ipAddress);
  }

  @override
  void didUpdateWidget(covariant _DasIpCard oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.das.ipAddress != widget.das.ipAddress) {
      _controller.text = widget.das.ipAddress;
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 2,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 16),
        child: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: Theme.of(context).primaryColor.withOpacity(0.1),
                shape: BoxShape.circle,
              ),
              child: Icon(Icons.router, color: Theme.of(context).primaryColor),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    widget.das.name,
                    style: const TextStyle(
                      fontWeight: FontWeight.bold,
                      fontSize: 16,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: _controller,
                    decoration: const InputDecoration(
                      labelText: 'IP Address',
                      isDense: true,
                      border: OutlineInputBorder(),
                    ),
                    keyboardType: TextInputType.number,
                  ),
                ],
              ),
            ),
            const SizedBox(width: 16),
            IconButton(
              icon: const Icon(Icons.save, color: Colors.blue),
              onPressed: () => widget.onSave(_controller.text),
              tooltip: 'Update IP',
            ),
          ],
        ),
      ),
    );
  }
}

// -----------------------------------------------------------------------------
// 3. Remarks Tab
// -----------------------------------------------------------------------------

class _RemarksTab extends StatefulWidget {
  const _RemarksTab();

  @override
  State<_RemarksTab> createState() => _RemarksTabState();
}

class _RemarksTabState extends State<_RemarksTab> {
  List<AcqRemark> _acquisitions = [];
  AcqRemark? _selectedAcq;
  late TextEditingController _remarkController;
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _remarkController = TextEditingController();
    _fetchRemarks();
  }

  @override
  void dispose() {
    _remarkController.dispose();
    super.dispose();
  }

  Future<void> _fetchRemarks() async {
    setState(() {
      _isLoading = true;
    });

    try {
      final globalState = context.read<GlobalState>();
      final channel = GrpcWebClientChannel.xhr(
        Uri.parse(globalState.serverUrl),
      );
      final client = CommunicationClient(channel);
      final request = ClientID()..id = globalState.clientId;

      final response = await client.getAcqRemarks(request);

      if (mounted) {
        setState(() {
          _acquisitions = response.acqRemarks;
          _isLoading = false;
        });
      }
    } catch (e) {
      debugPrint('Error fetching remarks: $e');
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('Failed to load remarks: $e')));
      }
    }
  }

  void _onAcquisitionChanged(AcqRemark? newValue) {
    setState(() {
      _selectedAcq = newValue;
      _remarkController.text = newValue?.remark ?? '';
    });
  }

  Future<void> _saveRemark() async {
    if (_selectedAcq == null) return;

    setState(() {
      _isLoading = true;
    });

    try {
      final globalState = context.read<GlobalState>();
      final channel = GrpcWebClientChannel.xhr(
        Uri.parse(globalState.serverUrl),
      );
      final client = CommunicationClient(channel);

      // Create a request object with updated remark
      final request = AcqRemark()
        ..date = _selectedAcq!.date
        ..time = _selectedAcq!.time
        ..remark = _remarkController.text;

      final response = await client.changeAcqRemark(request);

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              response.ok
                  ? response.message
                  : 'Failed to update remark: ${response.message}',
            ),
            backgroundColor: response.ok ? Colors.green : Colors.red,
          ),
        );

        if (response.ok) {
          // Update local state
          setState(() {
            _selectedAcq!.remark = _remarkController.text;
          });
        }
      }
    } catch (e) {
      debugPrint('Error saving remark: $e');
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e'), backgroundColor: Colors.red),
        );
      }
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading && _acquisitions.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_acquisitions.isEmpty) {
      return const Center(child: Text("No acquisitions found."));
    }

    return SingleChildScrollView(
      padding: const EdgeInsets.all(24.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Card(
            elevation: 4,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(12),
            ),
            child: Padding(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    'Select Acquisition',
                    style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
                  ),
                  const SizedBox(height: 16),
                  DropdownButtonFormField<AcqRemark>(
                    value: _selectedAcq,
                    decoration: const InputDecoration(
                      labelText: 'Acquisition Date/Time',
                      border: OutlineInputBorder(),
                      prefixIcon: Icon(Icons.calendar_today),
                    ),
                    items: _acquisitions.map((acq) {
                      return DropdownMenuItem(
                        value: acq,
                        child: Text('${acq.date} ${acq.time}'),
                      );
                    }).toList(),
                    onChanged: _onAcquisitionChanged,
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 24),
          if (_selectedAcq != null)
            Card(
              elevation: 4,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
              child: Padding(
                padding: const EdgeInsets.all(16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Edit Remark',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 16),
                    TextField(
                      controller: _remarkController,
                      maxLines: 5,
                      decoration: const InputDecoration(
                        labelText: 'Remarks',
                        alignLabelWithHint: true,
                        border: OutlineInputBorder(),
                        hintText: 'Enter observation details here...',
                      ),
                    ),
                    const SizedBox(height: 16),
                    Align(
                      alignment: Alignment.centerRight,
                      child: ElevatedButton.icon(
                        onPressed: _isLoading ? null : _saveRemark,
                        icon: const Icon(Icons.save),
                        label: _isLoading
                            ? const SizedBox(
                                width: 16,
                                height: 16,
                                child: CircularProgressIndicator(
                                  strokeWidth: 2,
                                  color: Colors.white,
                                ),
                              )
                            : const Text('Save Remark'),
                        style: ElevatedButton.styleFrom(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 32,
                            vertical: 12,
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            )
          else
            const Center(
              child: Padding(
                padding: EdgeInsets.only(top: 48.0),
                child: Text(
                  'Please select an acquisition to view or edit remarks.',
                  style: TextStyle(color: Colors.grey, fontSize: 16),
                ),
              ),
            ),
        ],
      ),
    );
  }
}
