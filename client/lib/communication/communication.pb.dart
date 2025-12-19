// This is a generated file - do not edit.
//
// Generated from communication.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class ClientID extends $pb.GeneratedMessage {
  factory ClientID({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  ClientID._();

  factory ClientID.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ClientID.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ClientID',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ClientID clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ClientID copyWith(void Function(ClientID) updates) =>
      super.copyWith((message) => updates(message as ClientID)) as ClientID;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ClientID create() => ClientID._();
  @$core.override
  ClientID createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ClientID getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ClientID>(create);
  static ClientID? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);
}

class Ack extends $pb.GeneratedMessage {
  factory Ack({
    $core.bool? ok,
    $core.String? message,
  }) {
    final result = create();
    if (ok != null) result.ok = ok;
    if (message != null) result.message = message;
    return result;
  }

  Ack._();

  factory Ack.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Ack.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Ack',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'ok')
    ..aOS(2, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Ack clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Ack copyWith(void Function(Ack) updates) =>
      super.copyWith((message) => updates(message as Ack)) as Ack;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Ack create() => Ack._();
  @$core.override
  Ack createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Ack getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Ack>(create);
  static Ack? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get ok => $_getBF(0);
  @$pb.TagNumber(1)
  set ok($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasOk() => $_has(0);
  @$pb.TagNumber(1)
  void clearOk() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get message => $_getSZ(1);
  @$pb.TagNumber(2)
  set message($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(2)
  void clearMessage() => $_clearField(2);
}

class ServerStatus extends $pb.GeneratedMessage {
  factory ServerStatus({
    $core.String? timestamp,
    $core.double? memory,
    $core.double? cpu,
  }) {
    final result = create();
    if (timestamp != null) result.timestamp = timestamp;
    if (memory != null) result.memory = memory;
    if (cpu != null) result.cpu = cpu;
    return result;
  }

  ServerStatus._();

  factory ServerStatus.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ServerStatus.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ServerStatus',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'timestamp')
    ..aD(2, _omitFieldNames ? '' : 'memory')
    ..aD(3, _omitFieldNames ? '' : 'cpu')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ServerStatus clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ServerStatus copyWith(void Function(ServerStatus) updates) =>
      super.copyWith((message) => updates(message as ServerStatus))
          as ServerStatus;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ServerStatus create() => ServerStatus._();
  @$core.override
  ServerStatus createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ServerStatus getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ServerStatus>(create);
  static ServerStatus? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get timestamp => $_getSZ(0);
  @$pb.TagNumber(1)
  set timestamp($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasTimestamp() => $_has(0);
  @$pb.TagNumber(1)
  void clearTimestamp() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.double get memory => $_getN(1);
  @$pb.TagNumber(2)
  set memory($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(2)
  $core.bool hasMemory() => $_has(1);
  @$pb.TagNumber(2)
  void clearMemory() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.double get cpu => $_getN(2);
  @$pb.TagNumber(3)
  set cpu($core.double value) => $_setDouble(2, value);
  @$pb.TagNumber(3)
  $core.bool hasCpu() => $_has(2);
  @$pb.TagNumber(3)
  void clearCpu() => $_clearField(3);
}

class IsWhitelistedResponse extends $pb.GeneratedMessage {
  factory IsWhitelistedResponse({
    $core.bool? whitelisted,
  }) {
    final result = create();
    if (whitelisted != null) result.whitelisted = whitelisted;
    return result;
  }

  IsWhitelistedResponse._();

  factory IsWhitelistedResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory IsWhitelistedResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'IsWhitelistedResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'whitelisted')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IsWhitelistedResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  IsWhitelistedResponse copyWith(
          void Function(IsWhitelistedResponse) updates) =>
      super.copyWith((message) => updates(message as IsWhitelistedResponse))
          as IsWhitelistedResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static IsWhitelistedResponse create() => IsWhitelistedResponse._();
  @$core.override
  IsWhitelistedResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static IsWhitelistedResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<IsWhitelistedResponse>(create);
  static IsWhitelistedResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get whitelisted => $_getBF(0);
  @$pb.TagNumber(1)
  set whitelisted($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasWhitelisted() => $_has(0);
  @$pb.TagNumber(1)
  void clearWhitelisted() => $_clearField(1);
}

class ServerDetails extends $pb.GeneratedMessage {
  factory ServerDetails({
    $core.String? satelliteName,
    $core.String? testPhase,
  }) {
    final result = create();
    if (satelliteName != null) result.satelliteName = satelliteName;
    if (testPhase != null) result.testPhase = testPhase;
    return result;
  }

  ServerDetails._();

  factory ServerDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ServerDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ServerDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'satelliteName', protoName: 'satelliteName')
    ..aOS(2, _omitFieldNames ? '' : 'testPhase', protoName: 'testPhase')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ServerDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ServerDetails copyWith(void Function(ServerDetails) updates) =>
      super.copyWith((message) => updates(message as ServerDetails))
          as ServerDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ServerDetails create() => ServerDetails._();
  @$core.override
  ServerDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ServerDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ServerDetails>(create);
  static ServerDetails? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get satelliteName => $_getSZ(0);
  @$pb.TagNumber(1)
  set satelliteName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSatelliteName() => $_has(0);
  @$pb.TagNumber(1)
  void clearSatelliteName() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get testPhase => $_getSZ(1);
  @$pb.TagNumber(2)
  set testPhase($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasTestPhase() => $_has(1);
  @$pb.TagNumber(2)
  void clearTestPhase() => $_clearField(2);
}

class LoginRequest extends $pb.GeneratedMessage {
  factory LoginRequest({
    $core.String? username,
    $core.String? password,
  }) {
    final result = create();
    if (username != null) result.username = username;
    if (password != null) result.password = password;
    return result;
  }

  LoginRequest._();

  factory LoginRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LoginRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LoginRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'username')
    ..aOS(2, _omitFieldNames ? '' : 'password')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LoginRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LoginRequest copyWith(void Function(LoginRequest) updates) =>
      super.copyWith((message) => updates(message as LoginRequest))
          as LoginRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LoginRequest create() => LoginRequest._();
  @$core.override
  LoginRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LoginRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LoginRequest>(create);
  static LoginRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get username => $_getSZ(0);
  @$pb.TagNumber(1)
  set username($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasUsername() => $_has(0);
  @$pb.TagNumber(1)
  void clearUsername() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get password => $_getSZ(1);
  @$pb.TagNumber(2)
  set password($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPassword() => $_has(1);
  @$pb.TagNumber(2)
  void clearPassword() => $_clearField(2);
}

class LoginResponse extends $pb.GeneratedMessage {
  factory LoginResponse({
    $core.bool? success,
    $core.Iterable<$core.String>? permissions,
  }) {
    final result = create();
    if (success != null) result.success = success;
    if (permissions != null) result.permissions.addAll(permissions);
    return result;
  }

  LoginResponse._();

  factory LoginResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LoginResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LoginResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'success')
    ..pPS(2, _omitFieldNames ? '' : 'permissions')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LoginResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LoginResponse copyWith(void Function(LoginResponse) updates) =>
      super.copyWith((message) => updates(message as LoginResponse))
          as LoginResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LoginResponse create() => LoginResponse._();
  @$core.override
  LoginResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LoginResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LoginResponse>(create);
  static LoginResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get success => $_getBF(0);
  @$pb.TagNumber(1)
  set success($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSuccess() => $_has(0);
  @$pb.TagNumber(1)
  void clearSuccess() => $_clearField(1);

  @$pb.TagNumber(2)
  $pb.PbList<$core.String> get permissions => $_getList(1);
}

class AcqDASMap extends $pb.GeneratedMessage {
  factory AcqDASMap({
    $core.String? acqMode,
    $core.Iterable<AcqDasDetails>? dasDetails,
  }) {
    final result = create();
    if (acqMode != null) result.acqMode = acqMode;
    if (dasDetails != null) result.dasDetails.addAll(dasDetails);
    return result;
  }

  AcqDASMap._();

  factory AcqDASMap.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AcqDASMap.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AcqDASMap',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'acqMode', protoName: 'acqMode')
    ..pPM<AcqDasDetails>(2, _omitFieldNames ? '' : 'dasDetails',
        protoName: 'dasDetails', subBuilder: AcqDasDetails.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AcqDASMap clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AcqDASMap copyWith(void Function(AcqDASMap) updates) =>
      super.copyWith((message) => updates(message as AcqDASMap)) as AcqDASMap;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AcqDASMap create() => AcqDASMap._();
  @$core.override
  AcqDASMap createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AcqDASMap getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AcqDASMap>(create);
  static AcqDASMap? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get acqMode => $_getSZ(0);
  @$pb.TagNumber(1)
  set acqMode($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAcqMode() => $_has(0);
  @$pb.TagNumber(1)
  void clearAcqMode() => $_clearField(1);

  @$pb.TagNumber(2)
  $pb.PbList<AcqDasDetails> get dasDetails => $_getList(1);
}

class AcqDasDetails extends $pb.GeneratedMessage {
  factory AcqDasDetails({
    $core.String? dasName,
    $core.int? dpuNumber,
  }) {
    final result = create();
    if (dasName != null) result.dasName = dasName;
    if (dpuNumber != null) result.dpuNumber = dpuNumber;
    return result;
  }

  AcqDasDetails._();

  factory AcqDasDetails.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AcqDasDetails.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AcqDasDetails',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'dasName', protoName: 'dasName')
    ..aI(2, _omitFieldNames ? '' : 'dpuNumber', protoName: 'dpuNumber')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AcqDasDetails clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AcqDasDetails copyWith(void Function(AcqDasDetails) updates) =>
      super.copyWith((message) => updates(message as AcqDasDetails))
          as AcqDasDetails;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AcqDasDetails create() => AcqDasDetails._();
  @$core.override
  AcqDasDetails createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AcqDasDetails getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AcqDasDetails>(create);
  static AcqDasDetails? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get dasName => $_getSZ(0);
  @$pb.TagNumber(1)
  set dasName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasDasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearDasName() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get dpuNumber => $_getIZ(1);
  @$pb.TagNumber(2)
  set dpuNumber($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasDpuNumber() => $_has(1);
  @$pb.TagNumber(2)
  void clearDpuNumber() => $_clearField(2);
}

class DASStatus extends $pb.GeneratedMessage {
  factory DASStatus({
    $core.Iterable<DASStatusResponse>? dasStatus,
  }) {
    final result = create();
    if (dasStatus != null) result.dasStatus.addAll(dasStatus);
    return result;
  }

  DASStatus._();

  factory DASStatus.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DASStatus.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DASStatus',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..pPM<DASStatusResponse>(1, _omitFieldNames ? '' : 'dasStatus',
        protoName: 'dasStatus', subBuilder: DASStatusResponse.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DASStatus clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DASStatus copyWith(void Function(DASStatus) updates) =>
      super.copyWith((message) => updates(message as DASStatus)) as DASStatus;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DASStatus create() => DASStatus._();
  @$core.override
  DASStatus createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DASStatus getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<DASStatus>(create);
  static DASStatus? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<DASStatusResponse> get dasStatus => $_getList(0);
}

class DASStatusResponse extends $pb.GeneratedMessage {
  factory DASStatusResponse({
    $core.String? dasName,
    $core.int? dpuNumber,
    $core.String? status,
    $core.bool? alarm,
  }) {
    final result = create();
    if (dasName != null) result.dasName = dasName;
    if (dpuNumber != null) result.dpuNumber = dpuNumber;
    if (status != null) result.status = status;
    if (alarm != null) result.alarm = alarm;
    return result;
  }

  DASStatusResponse._();

  factory DASStatusResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DASStatusResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DASStatusResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'dasName', protoName: 'dasName')
    ..aI(2, _omitFieldNames ? '' : 'dpuNumber', protoName: 'dpuNumber')
    ..aOS(3, _omitFieldNames ? '' : 'status')
    ..aOB(4, _omitFieldNames ? '' : 'alarm')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DASStatusResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DASStatusResponse copyWith(void Function(DASStatusResponse) updates) =>
      super.copyWith((message) => updates(message as DASStatusResponse))
          as DASStatusResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DASStatusResponse create() => DASStatusResponse._();
  @$core.override
  DASStatusResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DASStatusResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DASStatusResponse>(create);
  static DASStatusResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get dasName => $_getSZ(0);
  @$pb.TagNumber(1)
  set dasName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasDasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearDasName() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get dpuNumber => $_getIZ(1);
  @$pb.TagNumber(2)
  set dpuNumber($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasDpuNumber() => $_has(1);
  @$pb.TagNumber(2)
  void clearDpuNumber() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get status => $_getSZ(2);
  @$pb.TagNumber(3)
  set status($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasStatus() => $_has(2);
  @$pb.TagNumber(3)
  void clearStatus() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.bool get alarm => $_getBF(3);
  @$pb.TagNumber(4)
  set alarm($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(4)
  $core.bool hasAlarm() => $_has(3);
  @$pb.TagNumber(4)
  void clearAlarm() => $_clearField(4);
}

class DASStatusRequest extends $pb.GeneratedMessage {
  factory DASStatusRequest({
    $core.String? id,
    $core.String? acqMode,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (acqMode != null) result.acqMode = acqMode;
    return result;
  }

  DASStatusRequest._();

  factory DASStatusRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DASStatusRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DASStatusRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'acqMode', protoName: 'acqMode')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DASStatusRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DASStatusRequest copyWith(void Function(DASStatusRequest) updates) =>
      super.copyWith((message) => updates(message as DASStatusRequest))
          as DASStatusRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DASStatusRequest create() => DASStatusRequest._();
  @$core.override
  DASStatusRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DASStatusRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DASStatusRequest>(create);
  static DASStatusRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get acqMode => $_getSZ(1);
  @$pb.TagNumber(2)
  set acqMode($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasAcqMode() => $_has(1);
  @$pb.TagNumber(2)
  void clearAcqMode() => $_clearField(2);
}

class AcquisitionParameters extends $pb.GeneratedMessage {
  factory AcquisitionParameters({
    $core.Iterable<$core.String>? acqModes,
    $core.Iterable<$core.String>? payloads,
    $core.Iterable<$core.String>? configNames,
    $core.Iterable<$core.String>? acqTypes,
    $core.Iterable<$core.String>? resultProfiles,
    $core.Iterable<AcqDASMap>? dasMap,
  }) {
    final result = create();
    if (acqModes != null) result.acqModes.addAll(acqModes);
    if (payloads != null) result.payloads.addAll(payloads);
    if (configNames != null) result.configNames.addAll(configNames);
    if (acqTypes != null) result.acqTypes.addAll(acqTypes);
    if (resultProfiles != null) result.resultProfiles.addAll(resultProfiles);
    if (dasMap != null) result.dasMap.addAll(dasMap);
    return result;
  }

  AcquisitionParameters._();

  factory AcquisitionParameters.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AcquisitionParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AcquisitionParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..pPS(1, _omitFieldNames ? '' : 'acqModes', protoName: 'acqModes')
    ..pPS(2, _omitFieldNames ? '' : 'payloads')
    ..pPS(3, _omitFieldNames ? '' : 'configNames', protoName: 'configNames')
    ..pPS(4, _omitFieldNames ? '' : 'acqTypes', protoName: 'acqTypes')
    ..pPS(5, _omitFieldNames ? '' : 'resultProfiles',
        protoName: 'resultProfiles')
    ..pPM<AcqDASMap>(6, _omitFieldNames ? '' : 'dasMap',
        protoName: 'dasMap', subBuilder: AcqDASMap.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AcquisitionParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AcquisitionParameters copyWith(
          void Function(AcquisitionParameters) updates) =>
      super.copyWith((message) => updates(message as AcquisitionParameters))
          as AcquisitionParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AcquisitionParameters create() => AcquisitionParameters._();
  @$core.override
  AcquisitionParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AcquisitionParameters getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AcquisitionParameters>(create);
  static AcquisitionParameters? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<$core.String> get acqModes => $_getList(0);

  @$pb.TagNumber(2)
  $pb.PbList<$core.String> get payloads => $_getList(1);

  @$pb.TagNumber(3)
  $pb.PbList<$core.String> get configNames => $_getList(2);

  @$pb.TagNumber(4)
  $pb.PbList<$core.String> get acqTypes => $_getList(3);

  @$pb.TagNumber(5)
  $pb.PbList<$core.String> get resultProfiles => $_getList(4);

  @$pb.TagNumber(6)
  $pb.PbList<AcqDASMap> get dasMap => $_getList(5);
}

class FileAcquisitionParameters extends $pb.GeneratedMessage {
  factory FileAcquisitionParameters({
    $core.Iterable<$core.String>? frameTypes,
    $core.Iterable<$core.String>? acqModes,
    $core.Iterable<$core.String>? payloads,
    $core.Iterable<$core.String>? configNames,
    $core.Iterable<$core.String>? resultProfiles,
    $core.Iterable<FrameTypeMap>? frameTypeMap,
  }) {
    final result = create();
    if (frameTypes != null) result.frameTypes.addAll(frameTypes);
    if (acqModes != null) result.acqModes.addAll(acqModes);
    if (payloads != null) result.payloads.addAll(payloads);
    if (configNames != null) result.configNames.addAll(configNames);
    if (resultProfiles != null) result.resultProfiles.addAll(resultProfiles);
    if (frameTypeMap != null) result.frameTypeMap.addAll(frameTypeMap);
    return result;
  }

  FileAcquisitionParameters._();

  factory FileAcquisitionParameters.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FileAcquisitionParameters.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FileAcquisitionParameters',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..pPS(1, _omitFieldNames ? '' : 'frameTypes', protoName: 'frameTypes')
    ..pPS(2, _omitFieldNames ? '' : 'acqModes', protoName: 'acqModes')
    ..pPS(3, _omitFieldNames ? '' : 'payloads')
    ..pPS(4, _omitFieldNames ? '' : 'configNames', protoName: 'configNames')
    ..pPS(5, _omitFieldNames ? '' : 'resultProfiles',
        protoName: 'resultProfiles')
    ..pPM<FrameTypeMap>(6, _omitFieldNames ? '' : 'frameTypeMap',
        protoName: 'frameTypeMap', subBuilder: FrameTypeMap.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FileAcquisitionParameters clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FileAcquisitionParameters copyWith(
          void Function(FileAcquisitionParameters) updates) =>
      super.copyWith((message) => updates(message as FileAcquisitionParameters))
          as FileAcquisitionParameters;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FileAcquisitionParameters create() => FileAcquisitionParameters._();
  @$core.override
  FileAcquisitionParameters createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FileAcquisitionParameters getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FileAcquisitionParameters>(create);
  static FileAcquisitionParameters? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<$core.String> get frameTypes => $_getList(0);

  @$pb.TagNumber(2)
  $pb.PbList<$core.String> get acqModes => $_getList(1);

  @$pb.TagNumber(3)
  $pb.PbList<$core.String> get payloads => $_getList(2);

  @$pb.TagNumber(4)
  $pb.PbList<$core.String> get configNames => $_getList(3);

  @$pb.TagNumber(5)
  $pb.PbList<$core.String> get resultProfiles => $_getList(4);

  @$pb.TagNumber(6)
  $pb.PbList<FrameTypeMap> get frameTypeMap => $_getList(5);
}

class FrameTypeMap extends $pb.GeneratedMessage {
  factory FrameTypeMap({
    $core.String? frameType,
    $core.Iterable<$core.String>? frameIdentifiers,
  }) {
    final result = create();
    if (frameType != null) result.frameType = frameType;
    if (frameIdentifiers != null)
      result.frameIdentifiers.addAll(frameIdentifiers);
    return result;
  }

  FrameTypeMap._();

  factory FrameTypeMap.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory FrameTypeMap.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'FrameTypeMap',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'frameType', protoName: 'frameType')
    ..pPS(2, _omitFieldNames ? '' : 'frameIdentifiers',
        protoName: 'frameIdentifiers')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FrameTypeMap clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  FrameTypeMap copyWith(void Function(FrameTypeMap) updates) =>
      super.copyWith((message) => updates(message as FrameTypeMap))
          as FrameTypeMap;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static FrameTypeMap create() => FrameTypeMap._();
  @$core.override
  FrameTypeMap createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static FrameTypeMap getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<FrameTypeMap>(create);
  static FrameTypeMap? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get frameType => $_getSZ(0);
  @$pb.TagNumber(1)
  set frameType($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasFrameType() => $_has(0);
  @$pb.TagNumber(1)
  void clearFrameType() => $_clearField(1);

  @$pb.TagNumber(2)
  $pb.PbList<$core.String> get frameIdentifiers => $_getList(1);
}

class TestPhases extends $pb.GeneratedMessage {
  factory TestPhases({
    $core.Iterable<$core.String>? testPhases,
  }) {
    final result = create();
    if (testPhases != null) result.testPhases.addAll(testPhases);
    return result;
  }

  TestPhases._();

  factory TestPhases.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TestPhases.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TestPhases',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..pPS(1, _omitFieldNames ? '' : 'testPhases', protoName: 'testPhases')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestPhases clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestPhases copyWith(void Function(TestPhases) updates) =>
      super.copyWith((message) => updates(message as TestPhases)) as TestPhases;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TestPhases create() => TestPhases._();
  @$core.override
  TestPhases createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TestPhases getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TestPhases>(create);
  static TestPhases? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<$core.String> get testPhases => $_getList(0);
}

class TestPhaseRequest extends $pb.GeneratedMessage {
  factory TestPhaseRequest({
    $core.String? id,
    $core.String? testPhase,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (testPhase != null) result.testPhase = testPhase;
    return result;
  }

  TestPhaseRequest._();

  factory TestPhaseRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TestPhaseRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TestPhaseRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'testPhase', protoName: 'testPhase')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestPhaseRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TestPhaseRequest copyWith(void Function(TestPhaseRequest) updates) =>
      super.copyWith((message) => updates(message as TestPhaseRequest))
          as TestPhaseRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TestPhaseRequest create() => TestPhaseRequest._();
  @$core.override
  TestPhaseRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TestPhaseRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TestPhaseRequest>(create);
  static TestPhaseRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get testPhase => $_getSZ(1);
  @$pb.TagNumber(2)
  set testPhase($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasTestPhase() => $_has(1);
  @$pb.TagNumber(2)
  void clearTestPhase() => $_clearField(2);
}

class DASIPAddresses extends $pb.GeneratedMessage {
  factory DASIPAddresses({
    $core.Iterable<DASIPAddress>? dasIPAddresses,
  }) {
    final result = create();
    if (dasIPAddresses != null) result.dasIPAddresses.addAll(dasIPAddresses);
    return result;
  }

  DASIPAddresses._();

  factory DASIPAddresses.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DASIPAddresses.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DASIPAddresses',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..pPM<DASIPAddress>(1, _omitFieldNames ? '' : 'dasIPAddresses',
        protoName: 'dasIPAddresses', subBuilder: DASIPAddress.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DASIPAddresses clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DASIPAddresses copyWith(void Function(DASIPAddresses) updates) =>
      super.copyWith((message) => updates(message as DASIPAddresses))
          as DASIPAddresses;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DASIPAddresses create() => DASIPAddresses._();
  @$core.override
  DASIPAddresses createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DASIPAddresses getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DASIPAddresses>(create);
  static DASIPAddresses? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<DASIPAddress> get dasIPAddresses => $_getList(0);
}

class DASIPAddress extends $pb.GeneratedMessage {
  factory DASIPAddress({
    $core.String? name,
    $core.String? ipAddress,
  }) {
    final result = create();
    if (name != null) result.name = name;
    if (ipAddress != null) result.ipAddress = ipAddress;
    return result;
  }

  DASIPAddress._();

  factory DASIPAddress.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DASIPAddress.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DASIPAddress',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'name')
    ..aOS(2, _omitFieldNames ? '' : 'ipAddress', protoName: 'ipAddress')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DASIPAddress clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DASIPAddress copyWith(void Function(DASIPAddress) updates) =>
      super.copyWith((message) => updates(message as DASIPAddress))
          as DASIPAddress;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DASIPAddress create() => DASIPAddress._();
  @$core.override
  DASIPAddress createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DASIPAddress getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DASIPAddress>(create);
  static DASIPAddress? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(1)
  set name($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearName() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get ipAddress => $_getSZ(1);
  @$pb.TagNumber(2)
  set ipAddress($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasIpAddress() => $_has(1);
  @$pb.TagNumber(2)
  void clearIpAddress() => $_clearField(2);
}

class AcqRemark extends $pb.GeneratedMessage {
  factory AcqRemark({
    $core.String? date,
    $core.String? time,
    $core.String? remark,
  }) {
    final result = create();
    if (date != null) result.date = date;
    if (time != null) result.time = time;
    if (remark != null) result.remark = remark;
    return result;
  }

  AcqRemark._();

  factory AcqRemark.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AcqRemark.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AcqRemark',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'date')
    ..aOS(2, _omitFieldNames ? '' : 'time')
    ..aOS(3, _omitFieldNames ? '' : 'remark')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AcqRemark clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AcqRemark copyWith(void Function(AcqRemark) updates) =>
      super.copyWith((message) => updates(message as AcqRemark)) as AcqRemark;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AcqRemark create() => AcqRemark._();
  @$core.override
  AcqRemark createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AcqRemark getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AcqRemark>(create);
  static AcqRemark? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get date => $_getSZ(0);
  @$pb.TagNumber(1)
  set date($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasDate() => $_has(0);
  @$pb.TagNumber(1)
  void clearDate() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get time => $_getSZ(1);
  @$pb.TagNumber(2)
  set time($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasTime() => $_has(1);
  @$pb.TagNumber(2)
  void clearTime() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get remark => $_getSZ(2);
  @$pb.TagNumber(3)
  set remark($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasRemark() => $_has(2);
  @$pb.TagNumber(3)
  void clearRemark() => $_clearField(3);
}

class AcqRemarks extends $pb.GeneratedMessage {
  factory AcqRemarks({
    $core.Iterable<AcqRemark>? acqRemarks,
  }) {
    final result = create();
    if (acqRemarks != null) result.acqRemarks.addAll(acqRemarks);
    return result;
  }

  AcqRemarks._();

  factory AcqRemarks.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AcqRemarks.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AcqRemarks',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'Communication'),
      createEmptyInstance: create)
    ..pPM<AcqRemark>(1, _omitFieldNames ? '' : 'acqRemarks',
        protoName: 'acqRemarks', subBuilder: AcqRemark.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AcqRemarks clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AcqRemarks copyWith(void Function(AcqRemarks) updates) =>
      super.copyWith((message) => updates(message as AcqRemarks)) as AcqRemarks;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AcqRemarks create() => AcqRemarks._();
  @$core.override
  AcqRemarks createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AcqRemarks getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AcqRemarks>(create);
  static AcqRemarks? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<AcqRemark> get acqRemarks => $_getList(0);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
