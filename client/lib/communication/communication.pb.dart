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

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
