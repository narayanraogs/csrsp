// This is a generated file - do not edit.
//
// Generated from communication.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import 'communication.pb.dart' as $0;

export 'communication.pb.dart';

@$pb.GrpcServiceName('Communication.Communication')
class CommunicationClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  CommunicationClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseStream<$0.ServerStatus> getServerStatus(
    $0.ClientID request, {
    $grpc.CallOptions? options,
  }) {
    return $createStreamingCall(
        _$getServerStatus, $async.Stream.fromIterable([request]),
        options: options);
  }

  $grpc.ResponseFuture<$0.IsWhitelistedResponse> isWhitelisted(
    $0.ClientID request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$isWhitelisted, request, options: options);
  }

  $grpc.ResponseFuture<$0.ServerDetails> getServerDetails(
    $0.ClientID request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getServerDetails, request, options: options);
  }

  // method descriptors

  static final _$getServerStatus =
      $grpc.ClientMethod<$0.ClientID, $0.ServerStatus>(
          '/Communication.Communication/GetServerStatus',
          ($0.ClientID value) => value.writeToBuffer(),
          $0.ServerStatus.fromBuffer);
  static final _$isWhitelisted =
      $grpc.ClientMethod<$0.ClientID, $0.IsWhitelistedResponse>(
          '/Communication.Communication/IsWhitelisted',
          ($0.ClientID value) => value.writeToBuffer(),
          $0.IsWhitelistedResponse.fromBuffer);
  static final _$getServerDetails =
      $grpc.ClientMethod<$0.ClientID, $0.ServerDetails>(
          '/Communication.Communication/GetServerDetails',
          ($0.ClientID value) => value.writeToBuffer(),
          $0.ServerDetails.fromBuffer);
}

@$pb.GrpcServiceName('Communication.Communication')
abstract class CommunicationServiceBase extends $grpc.Service {
  $core.String get $name => 'Communication.Communication';

  CommunicationServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.ClientID, $0.ServerStatus>(
        'GetServerStatus',
        getServerStatus_Pre,
        false,
        true,
        ($core.List<$core.int> value) => $0.ClientID.fromBuffer(value),
        ($0.ServerStatus value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ClientID, $0.IsWhitelistedResponse>(
        'IsWhitelisted',
        isWhitelisted_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ClientID.fromBuffer(value),
        ($0.IsWhitelistedResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ClientID, $0.ServerDetails>(
        'GetServerDetails',
        getServerDetails_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ClientID.fromBuffer(value),
        ($0.ServerDetails value) => value.writeToBuffer()));
  }

  $async.Stream<$0.ServerStatus> getServerStatus_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.ClientID> $request) async* {
    yield* getServerStatus($call, await $request);
  }

  $async.Stream<$0.ServerStatus> getServerStatus(
      $grpc.ServiceCall call, $0.ClientID request);

  $async.Future<$0.IsWhitelistedResponse> isWhitelisted_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.ClientID> $request) async {
    return isWhitelisted($call, await $request);
  }

  $async.Future<$0.IsWhitelistedResponse> isWhitelisted(
      $grpc.ServiceCall call, $0.ClientID request);

  $async.Future<$0.ServerDetails> getServerDetails_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.ClientID> $request) async {
    return getServerDetails($call, await $request);
  }

  $async.Future<$0.ServerDetails> getServerDetails(
      $grpc.ServiceCall call, $0.ClientID request);
}
