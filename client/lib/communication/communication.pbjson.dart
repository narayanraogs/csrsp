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
