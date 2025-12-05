// This is a generated file - do not edit.
//
// Generated from communication.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports
// ignore_for_file: unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use clientIDDescriptor instead')
const ClientID$json = {
  '1': 'ClientID',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `ClientID`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List clientIDDescriptor =
    $convert.base64Decode('CghDbGllbnRJRBIOCgJpZBgBIAEoCVICaWQ=');

@$core.Deprecated('Use serverStatusDescriptor instead')
const ServerStatus$json = {
  '1': 'ServerStatus',
  '2': [
    {'1': 'timestamp', '3': 1, '4': 1, '5': 9, '10': 'timestamp'},
    {'1': 'memory', '3': 2, '4': 1, '5': 1, '10': 'memory'},
    {'1': 'cpu', '3': 3, '4': 1, '5': 1, '10': 'cpu'},
  ],
};

/// Descriptor for `ServerStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serverStatusDescriptor = $convert.base64Decode(
    'CgxTZXJ2ZXJTdGF0dXMSHAoJdGltZXN0YW1wGAEgASgJUgl0aW1lc3RhbXASFgoGbWVtb3J5GA'
    'IgASgBUgZtZW1vcnkSEAoDY3B1GAMgASgBUgNjcHU=');

@$core.Deprecated('Use isWhitelistedResponseDescriptor instead')
const IsWhitelistedResponse$json = {
  '1': 'IsWhitelistedResponse',
  '2': [
    {'1': 'whitelisted', '3': 1, '4': 1, '5': 8, '10': 'whitelisted'},
  ],
};

/// Descriptor for `IsWhitelistedResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List isWhitelistedResponseDescriptor = $convert.base64Decode(
    'ChVJc1doaXRlbGlzdGVkUmVzcG9uc2USIAoLd2hpdGVsaXN0ZWQYASABKAhSC3doaXRlbGlzdG'
    'Vk');

@$core.Deprecated('Use serverDetailsDescriptor instead')
const ServerDetails$json = {
  '1': 'ServerDetails',
  '2': [
    {'1': 'satelliteName', '3': 1, '4': 1, '5': 9, '10': 'satelliteName'},
    {'1': 'testPhase', '3': 2, '4': 1, '5': 9, '10': 'testPhase'},
  ],
};

/// Descriptor for `ServerDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serverDetailsDescriptor = $convert.base64Decode(
    'Cg1TZXJ2ZXJEZXRhaWxzEiQKDXNhdGVsbGl0ZU5hbWUYASABKAlSDXNhdGVsbGl0ZU5hbWUSHA'
    'oJdGVzdFBoYXNlGAIgASgJUgl0ZXN0UGhhc2U=');

@$core.Deprecated('Use loginRequestDescriptor instead')
const LoginRequest$json = {
  '1': 'LoginRequest',
  '2': [
    {'1': 'username', '3': 1, '4': 1, '5': 9, '10': 'username'},
    {'1': 'password', '3': 2, '4': 1, '5': 9, '10': 'password'},
  ],
};

/// Descriptor for `LoginRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loginRequestDescriptor = $convert.base64Decode(
    'CgxMb2dpblJlcXVlc3QSGgoIdXNlcm5hbWUYASABKAlSCHVzZXJuYW1lEhoKCHBhc3N3b3JkGA'
    'IgASgJUghwYXNzd29yZA==');

@$core.Deprecated('Use loginResponseDescriptor instead')
const LoginResponse$json = {
  '1': 'LoginResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'permissions', '3': 2, '4': 3, '5': 9, '10': 'permissions'},
  ],
};

/// Descriptor for `LoginResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loginResponseDescriptor = $convert.base64Decode(
    'Cg1Mb2dpblJlc3BvbnNlEhgKB3N1Y2Nlc3MYASABKAhSB3N1Y2Nlc3MSIAoLcGVybWlzc2lvbn'
    'MYAiADKAlSC3Blcm1pc3Npb25z');

@$core.Deprecated('Use acqDASMapDescriptor instead')
const AcqDASMap$json = {
  '1': 'AcqDASMap',
  '2': [
    {'1': 'acqMode', '3': 1, '4': 1, '5': 9, '10': 'acqMode'},
    {
      '1': 'dasDetails',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.Communication.AcqDasDetails',
      '10': 'dasDetails'
    },
  ],
};

/// Descriptor for `AcqDASMap`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List acqDASMapDescriptor = $convert.base64Decode(
    'CglBY3FEQVNNYXASGAoHYWNxTW9kZRgBIAEoCVIHYWNxTW9kZRI8CgpkYXNEZXRhaWxzGAIgAy'
    'gLMhwuQ29tbXVuaWNhdGlvbi5BY3FEYXNEZXRhaWxzUgpkYXNEZXRhaWxz');

@$core.Deprecated('Use acqDasDetailsDescriptor instead')
const AcqDasDetails$json = {
  '1': 'AcqDasDetails',
  '2': [
    {'1': 'dasName', '3': 1, '4': 1, '5': 9, '10': 'dasName'},
    {'1': 'dpuNumber', '3': 2, '4': 1, '5': 5, '10': 'dpuNumber'},
  ],
};

/// Descriptor for `AcqDasDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List acqDasDetailsDescriptor = $convert.base64Decode(
    'Cg1BY3FEYXNEZXRhaWxzEhgKB2Rhc05hbWUYASABKAlSB2Rhc05hbWUSHAoJZHB1TnVtYmVyGA'
    'IgASgFUglkcHVOdW1iZXI=');

@$core.Deprecated('Use acquisitionParametersDescriptor instead')
const AcquisitionParameters$json = {
  '1': 'AcquisitionParameters',
  '2': [
    {'1': 'acqModes', '3': 1, '4': 3, '5': 9, '10': 'acqModes'},
    {'1': 'payloads', '3': 2, '4': 3, '5': 9, '10': 'payloads'},
    {'1': 'configNames', '3': 3, '4': 3, '5': 9, '10': 'configNames'},
    {'1': 'acqTypes', '3': 4, '4': 3, '5': 9, '10': 'acqTypes'},
    {'1': 'resultProfiles', '3': 5, '4': 3, '5': 9, '10': 'resultProfiles'},
    {
      '1': 'dasMap',
      '3': 6,
      '4': 3,
      '5': 11,
      '6': '.Communication.AcqDASMap',
      '10': 'dasMap'
    },
  ],
};

/// Descriptor for `AcquisitionParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List acquisitionParametersDescriptor = $convert.base64Decode(
    'ChVBY3F1aXNpdGlvblBhcmFtZXRlcnMSGgoIYWNxTW9kZXMYASADKAlSCGFjcU1vZGVzEhoKCH'
    'BheWxvYWRzGAIgAygJUghwYXlsb2FkcxIgCgtjb25maWdOYW1lcxgDIAMoCVILY29uZmlnTmFt'
    'ZXMSGgoIYWNxVHlwZXMYBCADKAlSCGFjcVR5cGVzEiYKDnJlc3VsdFByb2ZpbGVzGAUgAygJUg'
    '5yZXN1bHRQcm9maWxlcxIwCgZkYXNNYXAYBiADKAsyGC5Db21tdW5pY2F0aW9uLkFjcURBU01h'
    'cFIGZGFzTWFw');
