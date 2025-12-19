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

  $grpc.ResponseFuture<$0.LoginResponse> login(
    $0.LoginRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$login, request, options: options);
  }

  $grpc.ResponseFuture<$0.AcquisitionParameters> getAcquisitionParameters(
    $0.ClientID request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAcquisitionParameters, request,
        options: options);
  }

  $grpc.ResponseStream<$0.DASStatus> getDASStatus(
    $0.DASStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createStreamingCall(
        _$getDASStatus, $async.Stream.fromIterable([request]),
        options: options);
  }

  $grpc.ResponseFuture<$0.FileAcquisitionParameters>
      getFileAcquisitionParameters(
    $0.ClientID request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getFileAcquisitionParameters, request,
        options: options);
  }

  $grpc.ResponseFuture<$0.TestPhases> getAllTestPhases(
    $0.ClientID request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAllTestPhases, request, options: options);
  }

  $grpc.ResponseFuture<$0.Ack> addTestPhase(
    $0.TestPhaseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$addTestPhase, request, options: options);
  }

  $grpc.ResponseFuture<$0.Ack> selectTestPhase(
    $0.TestPhaseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$selectTestPhase, request, options: options);
  }

  $grpc.ResponseFuture<$0.DASIPAddresses> getDASIPAddresses(
    $0.ClientID request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDASIPAddresses, request, options: options);
  }

  $grpc.ResponseFuture<$0.Ack> changeDASIPAddress(
    $0.DASIPAddress request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$changeDASIPAddress, request, options: options);
  }

  $grpc.ResponseFuture<$0.AcqRemarks> getAcqRemarks(
    $0.ClientID request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAcqRemarks, request, options: options);
  }

  $grpc.ResponseFuture<$0.Ack> changeAcqRemark(
    $0.AcqRemark request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$changeAcqRemark, request, options: options);
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
  static final _$login = $grpc.ClientMethod<$0.LoginRequest, $0.LoginResponse>(
      '/Communication.Communication/Login',
      ($0.LoginRequest value) => value.writeToBuffer(),
      $0.LoginResponse.fromBuffer);
  static final _$getAcquisitionParameters =
      $grpc.ClientMethod<$0.ClientID, $0.AcquisitionParameters>(
          '/Communication.Communication/GetAcquisitionParameters',
          ($0.ClientID value) => value.writeToBuffer(),
          $0.AcquisitionParameters.fromBuffer);
  static final _$getDASStatus =
      $grpc.ClientMethod<$0.DASStatusRequest, $0.DASStatus>(
          '/Communication.Communication/GetDASStatus',
          ($0.DASStatusRequest value) => value.writeToBuffer(),
          $0.DASStatus.fromBuffer);
  static final _$getFileAcquisitionParameters =
      $grpc.ClientMethod<$0.ClientID, $0.FileAcquisitionParameters>(
          '/Communication.Communication/GetFileAcquisitionParameters',
          ($0.ClientID value) => value.writeToBuffer(),
          $0.FileAcquisitionParameters.fromBuffer);
  static final _$getAllTestPhases =
      $grpc.ClientMethod<$0.ClientID, $0.TestPhases>(
          '/Communication.Communication/GetAllTestPhases',
          ($0.ClientID value) => value.writeToBuffer(),
          $0.TestPhases.fromBuffer);
  static final _$addTestPhase = $grpc.ClientMethod<$0.TestPhaseRequest, $0.Ack>(
      '/Communication.Communication/AddTestPhase',
      ($0.TestPhaseRequest value) => value.writeToBuffer(),
      $0.Ack.fromBuffer);
  static final _$selectTestPhase =
      $grpc.ClientMethod<$0.TestPhaseRequest, $0.Ack>(
          '/Communication.Communication/SelectTestPhase',
          ($0.TestPhaseRequest value) => value.writeToBuffer(),
          $0.Ack.fromBuffer);
  static final _$getDASIPAddresses =
      $grpc.ClientMethod<$0.ClientID, $0.DASIPAddresses>(
          '/Communication.Communication/GetDASIPAddresses',
          ($0.ClientID value) => value.writeToBuffer(),
          $0.DASIPAddresses.fromBuffer);
  static final _$changeDASIPAddress =
      $grpc.ClientMethod<$0.DASIPAddress, $0.Ack>(
          '/Communication.Communication/ChangeDASIPAddress',
          ($0.DASIPAddress value) => value.writeToBuffer(),
          $0.Ack.fromBuffer);
  static final _$getAcqRemarks = $grpc.ClientMethod<$0.ClientID, $0.AcqRemarks>(
      '/Communication.Communication/GetAcqRemarks',
      ($0.ClientID value) => value.writeToBuffer(),
      $0.AcqRemarks.fromBuffer);
  static final _$changeAcqRemark = $grpc.ClientMethod<$0.AcqRemark, $0.Ack>(
      '/Communication.Communication/ChangeAcqRemark',
      ($0.AcqRemark value) => value.writeToBuffer(),
      $0.Ack.fromBuffer);
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
    $addMethod($grpc.ServiceMethod<$0.LoginRequest, $0.LoginResponse>(
        'Login',
        login_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.LoginRequest.fromBuffer(value),
        ($0.LoginResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ClientID, $0.AcquisitionParameters>(
        'GetAcquisitionParameters',
        getAcquisitionParameters_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ClientID.fromBuffer(value),
        ($0.AcquisitionParameters value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DASStatusRequest, $0.DASStatus>(
        'GetDASStatus',
        getDASStatus_Pre,
        false,
        true,
        ($core.List<$core.int> value) => $0.DASStatusRequest.fromBuffer(value),
        ($0.DASStatus value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ClientID, $0.FileAcquisitionParameters>(
        'GetFileAcquisitionParameters',
        getFileAcquisitionParameters_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ClientID.fromBuffer(value),
        ($0.FileAcquisitionParameters value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ClientID, $0.TestPhases>(
        'GetAllTestPhases',
        getAllTestPhases_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ClientID.fromBuffer(value),
        ($0.TestPhases value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TestPhaseRequest, $0.Ack>(
        'AddTestPhase',
        addTestPhase_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.TestPhaseRequest.fromBuffer(value),
        ($0.Ack value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TestPhaseRequest, $0.Ack>(
        'SelectTestPhase',
        selectTestPhase_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.TestPhaseRequest.fromBuffer(value),
        ($0.Ack value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ClientID, $0.DASIPAddresses>(
        'GetDASIPAddresses',
        getDASIPAddresses_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ClientID.fromBuffer(value),
        ($0.DASIPAddresses value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DASIPAddress, $0.Ack>(
        'ChangeDASIPAddress',
        changeDASIPAddress_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.DASIPAddress.fromBuffer(value),
        ($0.Ack value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ClientID, $0.AcqRemarks>(
        'GetAcqRemarks',
        getAcqRemarks_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ClientID.fromBuffer(value),
        ($0.AcqRemarks value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AcqRemark, $0.Ack>(
        'ChangeAcqRemark',
        changeAcqRemark_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.AcqRemark.fromBuffer(value),
        ($0.Ack value) => value.writeToBuffer()));
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

  $async.Future<$0.LoginResponse> login_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.LoginRequest> $request) async {
    return login($call, await $request);
  }

  $async.Future<$0.LoginResponse> login(
      $grpc.ServiceCall call, $0.LoginRequest request);

  $async.Future<$0.AcquisitionParameters> getAcquisitionParameters_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.ClientID> $request) async {
    return getAcquisitionParameters($call, await $request);
  }

  $async.Future<$0.AcquisitionParameters> getAcquisitionParameters(
      $grpc.ServiceCall call, $0.ClientID request);

  $async.Stream<$0.DASStatus> getDASStatus_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DASStatusRequest> $request) async* {
    yield* getDASStatus($call, await $request);
  }

  $async.Stream<$0.DASStatus> getDASStatus(
      $grpc.ServiceCall call, $0.DASStatusRequest request);

  $async.Future<$0.FileAcquisitionParameters> getFileAcquisitionParameters_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.ClientID> $request) async {
    return getFileAcquisitionParameters($call, await $request);
  }

  $async.Future<$0.FileAcquisitionParameters> getFileAcquisitionParameters(
      $grpc.ServiceCall call, $0.ClientID request);

  $async.Future<$0.TestPhases> getAllTestPhases_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.ClientID> $request) async {
    return getAllTestPhases($call, await $request);
  }

  $async.Future<$0.TestPhases> getAllTestPhases(
      $grpc.ServiceCall call, $0.ClientID request);

  $async.Future<$0.Ack> addTestPhase_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TestPhaseRequest> $request) async {
    return addTestPhase($call, await $request);
  }

  $async.Future<$0.Ack> addTestPhase(
      $grpc.ServiceCall call, $0.TestPhaseRequest request);

  $async.Future<$0.Ack> selectTestPhase_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TestPhaseRequest> $request) async {
    return selectTestPhase($call, await $request);
  }

  $async.Future<$0.Ack> selectTestPhase(
      $grpc.ServiceCall call, $0.TestPhaseRequest request);

  $async.Future<$0.DASIPAddresses> getDASIPAddresses_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.ClientID> $request) async {
    return getDASIPAddresses($call, await $request);
  }

  $async.Future<$0.DASIPAddresses> getDASIPAddresses(
      $grpc.ServiceCall call, $0.ClientID request);

  $async.Future<$0.Ack> changeDASIPAddress_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.DASIPAddress> $request) async {
    return changeDASIPAddress($call, await $request);
  }

  $async.Future<$0.Ack> changeDASIPAddress(
      $grpc.ServiceCall call, $0.DASIPAddress request);

  $async.Future<$0.AcqRemarks> getAcqRemarks_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.ClientID> $request) async {
    return getAcqRemarks($call, await $request);
  }

  $async.Future<$0.AcqRemarks> getAcqRemarks(
      $grpc.ServiceCall call, $0.ClientID request);

  $async.Future<$0.Ack> changeAcqRemark_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.AcqRemark> $request) async {
    return changeAcqRemark($call, await $request);
  }

  $async.Future<$0.Ack> changeAcqRemark(
      $grpc.ServiceCall call, $0.AcqRemark request);
}
