/*eslint-disable block-scoped-var, id-length, no-control-regex, no-magic-numbers, no-prototype-builtins, no-redeclare, no-shadow, no-var, sort-vars*/
import * as $protobuf from "protobufjs/minimal";
// Common aliases
const $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;
// Exported root namespace
const $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});
export const messagespb = $root.messagespb = (() => {
    /**
     * Namespace messagespb.
     * @exports messagespb
     * @namespace
     */
    const messagespb = {};
    messagespb.MessagesService = (function() {
        /**
         * Constructs a new MessagesService service.
         * @memberof messagespb
         * @classdesc Represents a MessagesService
         * @extends $protobuf.rpc.Service
         * @constructor
         * @param {$protobuf.RPCImpl} rpcImpl RPC implementation
         * @param {boolean} [requestDelimited=false] Whether requests are length-delimited
         * @param {boolean} [responseDelimited=false] Whether responses are length-delimited
         */
        function MessagesService(rpcImpl, requestDelimited, responseDelimited) {
            $protobuf.rpc.Service.call(this, rpcImpl, requestDelimited, responseDelimited);
        }
        (MessagesService.prototype = Object.create($protobuf.rpc.Service.prototype)).constructor = MessagesService;
        /**
         * Creates new MessagesService service using the specified rpc implementation.
         * @function create
         * @memberof messagespb.MessagesService
         * @static
         * @param {$protobuf.RPCImpl} rpcImpl RPC implementation
         * @param {boolean} [requestDelimited=false] Whether requests are length-delimited
         * @param {boolean} [responseDelimited=false] Whether responses are length-delimited
         * @returns {MessagesService} RPC service. Useful where requests and/or responses are streamed.
         */
        MessagesService.create = function create(rpcImpl, requestDelimited, responseDelimited) {
            return new this(rpcImpl, requestDelimited, responseDelimited);
        };
        /**
         * Callback as used by {@link messagespb.MessagesService#startConversation}.
         * @memberof messagespb.MessagesService
         * @typedef StartConversationCallback
         * @type {function}
         * @param {Error|null} error Error, if any
         * @param {messagespb.StartConversationResponse} [response] StartConversationResponse
         */
        /**
         * Calls StartConversation.
         * @function startConversation
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.IStartConversationRequest} request StartConversationRequest message or plain object
         * @param {messagespb.MessagesService.StartConversationCallback} callback Node-style callback called with the error, if any, and StartConversationResponse
         * @returns {undefined}
         * @variation 1
         */
        Object.defineProperty(MessagesService.prototype.startConversation = function startConversation(request, callback) {
            return this.rpcCall(startConversation, $root.messagespb.StartConversationRequest, $root.messagespb.StartConversationResponse, request, callback);
        }, "name", { value: "StartConversation" });
        /**
         * Calls StartConversation.
         * @function startConversation
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.IStartConversationRequest} request StartConversationRequest message or plain object
         * @returns {Promise<messagespb.StartConversationResponse>} Promise
         * @variation 2
         */
        /**
         * Callback as used by {@link messagespb.MessagesService#restoreConversation}.
         * @memberof messagespb.MessagesService
         * @typedef RestoreConversationCallback
         * @type {function}
         * @param {Error|null} error Error, if any
         * @param {messagespb.RestoreConversationResponse} [response] RestoreConversationResponse
         */
        /**
         * Calls RestoreConversation.
         * @function restoreConversation
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.IRestoreConversationRequest} request RestoreConversationRequest message or plain object
         * @param {messagespb.MessagesService.RestoreConversationCallback} callback Node-style callback called with the error, if any, and RestoreConversationResponse
         * @returns {undefined}
         * @variation 1
         */
        Object.defineProperty(MessagesService.prototype.restoreConversation = function restoreConversation(request, callback) {
            return this.rpcCall(restoreConversation, $root.messagespb.RestoreConversationRequest, $root.messagespb.RestoreConversationResponse, request, callback);
        }, "name", { value: "RestoreConversation" });
        /**
         * Calls RestoreConversation.
         * @function restoreConversation
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.IRestoreConversationRequest} request RestoreConversationRequest message or plain object
         * @returns {Promise<messagespb.RestoreConversationResponse>} Promise
         * @variation 2
         */
        /**
         * Callback as used by {@link messagespb.MessagesService#archiveConversation}.
         * @memberof messagespb.MessagesService
         * @typedef ArchiveConversationCallback
         * @type {function}
         * @param {Error|null} error Error, if any
         * @param {messagespb.ArchiveConversationResponse} [response] ArchiveConversationResponse
         */
        /**
         * Calls ArchiveConversation.
         * @function archiveConversation
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.IArchiveConversationRequest} request ArchiveConversationRequest message or plain object
         * @param {messagespb.MessagesService.ArchiveConversationCallback} callback Node-style callback called with the error, if any, and ArchiveConversationResponse
         * @returns {undefined}
         * @variation 1
         */
        Object.defineProperty(MessagesService.prototype.archiveConversation = function archiveConversation(request, callback) {
            return this.rpcCall(archiveConversation, $root.messagespb.ArchiveConversationRequest, $root.messagespb.ArchiveConversationResponse, request, callback);
        }, "name", { value: "ArchiveConversation" });
        /**
         * Calls ArchiveConversation.
         * @function archiveConversation
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.IArchiveConversationRequest} request ArchiveConversationRequest message or plain object
         * @returns {Promise<messagespb.ArchiveConversationResponse>} Promise
         * @variation 2
         */
        /**
         * Callback as used by {@link messagespb.MessagesService#getConversation}.
         * @memberof messagespb.MessagesService
         * @typedef GetConversationCallback
         * @type {function}
         * @param {Error|null} error Error, if any
         * @param {messagespb.GetConversationResponse} [response] GetConversationResponse
         */
        /**
         * Calls GetConversation.
         * @function getConversation
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.IGetConversationRequest} request GetConversationRequest message or plain object
         * @param {messagespb.MessagesService.GetConversationCallback} callback Node-style callback called with the error, if any, and GetConversationResponse
         * @returns {undefined}
         * @variation 1
         */
        Object.defineProperty(MessagesService.prototype.getConversation = function getConversation(request, callback) {
            return this.rpcCall(getConversation, $root.messagespb.GetConversationRequest, $root.messagespb.GetConversationResponse, request, callback);
        }, "name", { value: "GetConversation" });
        /**
         * Calls GetConversation.
         * @function getConversation
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.IGetConversationRequest} request GetConversationRequest message or plain object
         * @returns {Promise<messagespb.GetConversationResponse>} Promise
         * @variation 2
         */
        /**
         * Callback as used by {@link messagespb.MessagesService#getConversations}.
         * @memberof messagespb.MessagesService
         * @typedef GetConversationsCallback
         * @type {function}
         * @param {Error|null} error Error, if any
         * @param {messagespb.GetConversationsResponse} [response] GetConversationsResponse
         */
        /**
         * Calls GetConversations.
         * @function getConversations
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.IGetConversationsRequest} request GetConversationsRequest message or plain object
         * @param {messagespb.MessagesService.GetConversationsCallback} callback Node-style callback called with the error, if any, and GetConversationsResponse
         * @returns {undefined}
         * @variation 1
         */
        Object.defineProperty(MessagesService.prototype.getConversations = function getConversations(request, callback) {
            return this.rpcCall(getConversations, $root.messagespb.GetConversationsRequest, $root.messagespb.GetConversationsResponse, request, callback);
        }, "name", { value: "GetConversations" });
        /**
         * Calls GetConversations.
         * @function getConversations
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.IGetConversationsRequest} request GetConversationsRequest message or plain object
         * @returns {Promise<messagespb.GetConversationsResponse>} Promise
         * @variation 2
         */
        /**
         * Callback as used by {@link messagespb.MessagesService#getActiveConversations}.
         * @memberof messagespb.MessagesService
         * @typedef GetActiveConversationsCallback
         * @type {function}
         * @param {Error|null} error Error, if any
         * @param {messagespb.GetActiveConversationsResponse} [response] GetActiveConversationsResponse
         */
        /**
         * Calls GetActiveConversations.
         * @function getActiveConversations
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.IGetActiveConversationsRequest} request GetActiveConversationsRequest message or plain object
         * @param {messagespb.MessagesService.GetActiveConversationsCallback} callback Node-style callback called with the error, if any, and GetActiveConversationsResponse
         * @returns {undefined}
         * @variation 1
         */
        Object.defineProperty(MessagesService.prototype.getActiveConversations = function getActiveConversations(request, callback) {
            return this.rpcCall(getActiveConversations, $root.messagespb.GetActiveConversationsRequest, $root.messagespb.GetActiveConversationsResponse, request, callback);
        }, "name", { value: "GetActiveConversations" });
        /**
         * Calls GetActiveConversations.
         * @function getActiveConversations
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.IGetActiveConversationsRequest} request GetActiveConversationsRequest message or plain object
         * @returns {Promise<messagespb.GetActiveConversationsResponse>} Promise
         * @variation 2
         */
        /**
         * Callback as used by {@link messagespb.MessagesService#sendMessage}.
         * @memberof messagespb.MessagesService
         * @typedef SendMessageCallback
         * @type {function}
         * @param {Error|null} error Error, if any
         * @param {messagespb.SendMessageResponse} [response] SendMessageResponse
         */
        /**
         * Calls SendMessage.
         * @function sendMessage
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.ISendMessageRequest} request SendMessageRequest message or plain object
         * @param {messagespb.MessagesService.SendMessageCallback} callback Node-style callback called with the error, if any, and SendMessageResponse
         * @returns {undefined}
         * @variation 1
         */
        Object.defineProperty(MessagesService.prototype.sendMessage = function sendMessage(request, callback) {
            return this.rpcCall(sendMessage, $root.messagespb.SendMessageRequest, $root.messagespb.SendMessageResponse, request, callback);
        }, "name", { value: "SendMessage" });
        /**
         * Calls SendMessage.
         * @function sendMessage
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.ISendMessageRequest} request SendMessageRequest message or plain object
         * @returns {Promise<messagespb.SendMessageResponse>} Promise
         * @variation 2
         */
        /**
         * Callback as used by {@link messagespb.MessagesService#deleteMessage}.
         * @memberof messagespb.MessagesService
         * @typedef DeleteMessageCallback
         * @type {function}
         * @param {Error|null} error Error, if any
         * @param {messagespb.DeleteMessageResponse} [response] DeleteMessageResponse
         */
        /**
         * Calls DeleteMessage.
         * @function deleteMessage
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.IDeleteMessageRequest} request DeleteMessageRequest message or plain object
         * @param {messagespb.MessagesService.DeleteMessageCallback} callback Node-style callback called with the error, if any, and DeleteMessageResponse
         * @returns {undefined}
         * @variation 1
         */
        Object.defineProperty(MessagesService.prototype.deleteMessage = function deleteMessage(request, callback) {
            return this.rpcCall(deleteMessage, $root.messagespb.DeleteMessageRequest, $root.messagespb.DeleteMessageResponse, request, callback);
        }, "name", { value: "DeleteMessage" });
        /**
         * Calls DeleteMessage.
         * @function deleteMessage
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.IDeleteMessageRequest} request DeleteMessageRequest message or plain object
         * @returns {Promise<messagespb.DeleteMessageResponse>} Promise
         * @variation 2
         */
        /**
         * Callback as used by {@link messagespb.MessagesService#getMessage}.
         * @memberof messagespb.MessagesService
         * @typedef GetMessageCallback
         * @type {function}
         * @param {Error|null} error Error, if any
         * @param {messagespb.GetMessageResponse} [response] GetMessageResponse
         */
        /**
         * Calls GetMessage.
         * @function getMessage
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.IGetMessageRequest} request GetMessageRequest message or plain object
         * @param {messagespb.MessagesService.GetMessageCallback} callback Node-style callback called with the error, if any, and GetMessageResponse
         * @returns {undefined}
         * @variation 1
         */
        Object.defineProperty(MessagesService.prototype.getMessage = function getMessage(request, callback) {
            return this.rpcCall(getMessage, $root.messagespb.GetMessageRequest, $root.messagespb.GetMessageResponse, request, callback);
        }, "name", { value: "GetMessage" });
        /**
         * Calls GetMessage.
         * @function getMessage
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.IGetMessageRequest} request GetMessageRequest message or plain object
         * @returns {Promise<messagespb.GetMessageResponse>} Promise
         * @variation 2
         */
        /**
         * Callback as used by {@link messagespb.MessagesService#getMessages}.
         * @memberof messagespb.MessagesService
         * @typedef GetMessagesCallback
         * @type {function}
         * @param {Error|null} error Error, if any
         * @param {messagespb.GetMessagesResponse} [response] GetMessagesResponse
         */
        /**
         * Calls GetMessages.
         * @function getMessages
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.IGetMessagesRequest} request GetMessagesRequest message or plain object
         * @param {messagespb.MessagesService.GetMessagesCallback} callback Node-style callback called with the error, if any, and GetMessagesResponse
         * @returns {undefined}
         * @variation 1
         */
        Object.defineProperty(MessagesService.prototype.getMessages = function getMessages(request, callback) {
            return this.rpcCall(getMessages, $root.messagespb.GetMessagesRequest, $root.messagespb.GetMessagesResponse, request, callback);
        }, "name", { value: "GetMessages" });
        /**
         * Calls GetMessages.
         * @function getMessages
         * @memberof messagespb.MessagesService
         * @instance
         * @param {messagespb.IGetMessagesRequest} request GetMessagesRequest message or plain object
         * @returns {Promise<messagespb.GetMessagesResponse>} Promise
         * @variation 2
         */
        return MessagesService;
    })();
    messagespb.DeleteMessageResponse = (function() {
        /**
         * Properties of a DeleteMessageResponse.
         * @memberof messagespb
         * @interface IDeleteMessageResponse
         */
        /**
         * Constructs a new DeleteMessageResponse.
         * @memberof messagespb
         * @classdesc Represents a DeleteMessageResponse.
         * @implements IDeleteMessageResponse
         * @constructor
         * @param {messagespb.IDeleteMessageResponse=} [properties] Properties to set
         */
        function DeleteMessageResponse(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * Creates a new DeleteMessageResponse instance using the specified properties.
         * @function create
         * @memberof messagespb.DeleteMessageResponse
         * @static
         * @param {messagespb.IDeleteMessageResponse=} [properties] Properties to set
         * @returns {messagespb.DeleteMessageResponse} DeleteMessageResponse instance
         */
        DeleteMessageResponse.create = function create(properties) {
            return new DeleteMessageResponse(properties);
        };
        /**
         * Encodes the specified DeleteMessageResponse message. Does not implicitly {@link messagespb.DeleteMessageResponse.verify|verify} messages.
         * @function encode
         * @memberof messagespb.DeleteMessageResponse
         * @static
         * @param {messagespb.IDeleteMessageResponse} message DeleteMessageResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        DeleteMessageResponse.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            return writer;
        };
        /**
         * Encodes the specified DeleteMessageResponse message, length delimited. Does not implicitly {@link messagespb.DeleteMessageResponse.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.DeleteMessageResponse
         * @static
         * @param {messagespb.IDeleteMessageResponse} message DeleteMessageResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        DeleteMessageResponse.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a DeleteMessageResponse message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.DeleteMessageResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.DeleteMessageResponse} DeleteMessageResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        DeleteMessageResponse.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.DeleteMessageResponse();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a DeleteMessageResponse message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.DeleteMessageResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.DeleteMessageResponse} DeleteMessageResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        DeleteMessageResponse.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a DeleteMessageResponse message.
         * @function verify
         * @memberof messagespb.DeleteMessageResponse
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        DeleteMessageResponse.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            return null;
        };
        /**
         * Creates a DeleteMessageResponse message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.DeleteMessageResponse
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.DeleteMessageResponse} DeleteMessageResponse
         */
        DeleteMessageResponse.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.DeleteMessageResponse)
                return object;
            return new $root.messagespb.DeleteMessageResponse();
        };
        /**
         * Creates a plain object from a DeleteMessageResponse message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.DeleteMessageResponse
         * @static
         * @param {messagespb.DeleteMessageResponse} message DeleteMessageResponse
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        DeleteMessageResponse.toObject = function toObject() {
            return {};
        };
        /**
         * Converts this DeleteMessageResponse to JSON.
         * @function toJSON
         * @memberof messagespb.DeleteMessageResponse
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        DeleteMessageResponse.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for DeleteMessageResponse
         * @function getTypeUrl
         * @memberof messagespb.DeleteMessageResponse
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        DeleteMessageResponse.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.DeleteMessageResponse";
        };
        return DeleteMessageResponse;
    })();
    messagespb.Conversation = (function() {
        /**
         * Properties of a Conversation.
         * @memberof messagespb
         * @interface IConversation
         * @property {string|null} [id] Conversation id
         * @property {string|null} [senderId] Conversation senderId
         * @property {string|null} [recipientId] Conversation recipientId
         * @property {string|null} [itemId] Conversation itemId
         * @property {string|null} [conversationStatus] Conversation conversationStatus
         */
        /**
         * Constructs a new Conversation.
         * @memberof messagespb
         * @classdesc Represents a Conversation.
         * @implements IConversation
         * @constructor
         * @param {messagespb.IConversation=} [properties] Properties to set
         */
        function Conversation(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * Conversation id.
         * @member {string} id
         * @memberof messagespb.Conversation
         * @instance
         */
        Conversation.prototype.id = "";
        /**
         * Conversation senderId.
         * @member {string} senderId
         * @memberof messagespb.Conversation
         * @instance
         */
        Conversation.prototype.senderId = "";
        /**
         * Conversation recipientId.
         * @member {string} recipientId
         * @memberof messagespb.Conversation
         * @instance
         */
        Conversation.prototype.recipientId = "";
        /**
         * Conversation itemId.
         * @member {string} itemId
         * @memberof messagespb.Conversation
         * @instance
         */
        Conversation.prototype.itemId = "";
        /**
         * Conversation conversationStatus.
         * @member {string} conversationStatus
         * @memberof messagespb.Conversation
         * @instance
         */
        Conversation.prototype.conversationStatus = "";
        /**
         * Creates a new Conversation instance using the specified properties.
         * @function create
         * @memberof messagespb.Conversation
         * @static
         * @param {messagespb.IConversation=} [properties] Properties to set
         * @returns {messagespb.Conversation} Conversation instance
         */
        Conversation.create = function create(properties) {
            return new Conversation(properties);
        };
        /**
         * Encodes the specified Conversation message. Does not implicitly {@link messagespb.Conversation.verify|verify} messages.
         * @function encode
         * @memberof messagespb.Conversation
         * @static
         * @param {messagespb.IConversation} message Conversation message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        Conversation.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
            if (message.senderId != null && Object.hasOwnProperty.call(message, "senderId"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.senderId);
            if (message.recipientId != null && Object.hasOwnProperty.call(message, "recipientId"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.recipientId);
            if (message.itemId != null && Object.hasOwnProperty.call(message, "itemId"))
                writer.uint32(/* id 4, wireType 2 =*/34).string(message.itemId);
            if (message.conversationStatus != null && Object.hasOwnProperty.call(message, "conversationStatus"))
                writer.uint32(/* id 5, wireType 2 =*/42).string(message.conversationStatus);
            return writer;
        };
        /**
         * Encodes the specified Conversation message, length delimited. Does not implicitly {@link messagespb.Conversation.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.Conversation
         * @static
         * @param {messagespb.IConversation} message Conversation message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        Conversation.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a Conversation message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.Conversation
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.Conversation} Conversation
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        Conversation.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.Conversation();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.id = reader.string();
                        break;
                    }
                case 2: {
                        message.senderId = reader.string();
                        break;
                    }
                case 3: {
                        message.recipientId = reader.string();
                        break;
                    }
                case 4: {
                        message.itemId = reader.string();
                        break;
                    }
                case 5: {
                        message.conversationStatus = reader.string();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a Conversation message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.Conversation
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.Conversation} Conversation
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        Conversation.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a Conversation message.
         * @function verify
         * @memberof messagespb.Conversation
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        Conversation.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.id != null && message.hasOwnProperty("id"))
                if (!$util.isString(message.id))
                    return "id: string expected";
            if (message.senderId != null && message.hasOwnProperty("senderId"))
                if (!$util.isString(message.senderId))
                    return "senderId: string expected";
            if (message.recipientId != null && message.hasOwnProperty("recipientId"))
                if (!$util.isString(message.recipientId))
                    return "recipientId: string expected";
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                if (!$util.isString(message.itemId))
                    return "itemId: string expected";
            if (message.conversationStatus != null && message.hasOwnProperty("conversationStatus"))
                if (!$util.isString(message.conversationStatus))
                    return "conversationStatus: string expected";
            return null;
        };
        /**
         * Creates a Conversation message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.Conversation
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.Conversation} Conversation
         */
        Conversation.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.Conversation)
                return object;
            let message = new $root.messagespb.Conversation();
            if (object.id != null)
                message.id = String(object.id);
            if (object.senderId != null)
                message.senderId = String(object.senderId);
            if (object.recipientId != null)
                message.recipientId = String(object.recipientId);
            if (object.itemId != null)
                message.itemId = String(object.itemId);
            if (object.conversationStatus != null)
                message.conversationStatus = String(object.conversationStatus);
            return message;
        };
        /**
         * Creates a plain object from a Conversation message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.Conversation
         * @static
         * @param {messagespb.Conversation} message Conversation
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        Conversation.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.id = "";
                object.senderId = "";
                object.recipientId = "";
                object.itemId = "";
                object.conversationStatus = "";
            }
            if (message.id != null && message.hasOwnProperty("id"))
                object.id = message.id;
            if (message.senderId != null && message.hasOwnProperty("senderId"))
                object.senderId = message.senderId;
            if (message.recipientId != null && message.hasOwnProperty("recipientId"))
                object.recipientId = message.recipientId;
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                object.itemId = message.itemId;
            if (message.conversationStatus != null && message.hasOwnProperty("conversationStatus"))
                object.conversationStatus = message.conversationStatus;
            return object;
        };
        /**
         * Converts this Conversation to JSON.
         * @function toJSON
         * @memberof messagespb.Conversation
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        Conversation.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for Conversation
         * @function getTypeUrl
         * @memberof messagespb.Conversation
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        Conversation.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.Conversation";
        };
        return Conversation;
    })();
    messagespb.Message = (function() {
        /**
         * Properties of a Message.
         * @memberof messagespb
         * @interface IMessage
         * @property {string|null} [id] Message id
         * @property {string|null} [conversationId] Message conversationId
         * @property {string|null} [senderId] Message senderId
         * @property {string|null} [recipientId] Message recipientId
         * @property {string|null} [itemId] Message itemId
         * @property {string|null} [body] Message body
         */
        /**
         * Constructs a new Message.
         * @memberof messagespb
         * @classdesc Represents a Message.
         * @implements IMessage
         * @constructor
         * @param {messagespb.IMessage=} [properties] Properties to set
         */
        function Message(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * Message id.
         * @member {string} id
         * @memberof messagespb.Message
         * @instance
         */
        Message.prototype.id = "";
        /**
         * Message conversationId.
         * @member {string} conversationId
         * @memberof messagespb.Message
         * @instance
         */
        Message.prototype.conversationId = "";
        /**
         * Message senderId.
         * @member {string} senderId
         * @memberof messagespb.Message
         * @instance
         */
        Message.prototype.senderId = "";
        /**
         * Message recipientId.
         * @member {string} recipientId
         * @memberof messagespb.Message
         * @instance
         */
        Message.prototype.recipientId = "";
        /**
         * Message itemId.
         * @member {string} itemId
         * @memberof messagespb.Message
         * @instance
         */
        Message.prototype.itemId = "";
        /**
         * Message body.
         * @member {string} body
         * @memberof messagespb.Message
         * @instance
         */
        Message.prototype.body = "";
        /**
         * Creates a new Message instance using the specified properties.
         * @function create
         * @memberof messagespb.Message
         * @static
         * @param {messagespb.IMessage=} [properties] Properties to set
         * @returns {messagespb.Message} Message instance
         */
        Message.create = function create(properties) {
            return new Message(properties);
        };
        /**
         * Encodes the specified Message message. Does not implicitly {@link messagespb.Message.verify|verify} messages.
         * @function encode
         * @memberof messagespb.Message
         * @static
         * @param {messagespb.IMessage} message Message message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        Message.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
            if (message.conversationId != null && Object.hasOwnProperty.call(message, "conversationId"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.conversationId);
            if (message.senderId != null && Object.hasOwnProperty.call(message, "senderId"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.senderId);
            if (message.recipientId != null && Object.hasOwnProperty.call(message, "recipientId"))
                writer.uint32(/* id 4, wireType 2 =*/34).string(message.recipientId);
            if (message.itemId != null && Object.hasOwnProperty.call(message, "itemId"))
                writer.uint32(/* id 5, wireType 2 =*/42).string(message.itemId);
            if (message.body != null && Object.hasOwnProperty.call(message, "body"))
                writer.uint32(/* id 6, wireType 2 =*/50).string(message.body);
            return writer;
        };
        /**
         * Encodes the specified Message message, length delimited. Does not implicitly {@link messagespb.Message.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.Message
         * @static
         * @param {messagespb.IMessage} message Message message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        Message.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a Message message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.Message
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.Message} Message
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        Message.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.Message();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.id = reader.string();
                        break;
                    }
                case 2: {
                        message.conversationId = reader.string();
                        break;
                    }
                case 3: {
                        message.senderId = reader.string();
                        break;
                    }
                case 4: {
                        message.recipientId = reader.string();
                        break;
                    }
                case 5: {
                        message.itemId = reader.string();
                        break;
                    }
                case 6: {
                        message.body = reader.string();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a Message message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.Message
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.Message} Message
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        Message.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a Message message.
         * @function verify
         * @memberof messagespb.Message
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        Message.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.id != null && message.hasOwnProperty("id"))
                if (!$util.isString(message.id))
                    return "id: string expected";
            if (message.conversationId != null && message.hasOwnProperty("conversationId"))
                if (!$util.isString(message.conversationId))
                    return "conversationId: string expected";
            if (message.senderId != null && message.hasOwnProperty("senderId"))
                if (!$util.isString(message.senderId))
                    return "senderId: string expected";
            if (message.recipientId != null && message.hasOwnProperty("recipientId"))
                if (!$util.isString(message.recipientId))
                    return "recipientId: string expected";
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                if (!$util.isString(message.itemId))
                    return "itemId: string expected";
            if (message.body != null && message.hasOwnProperty("body"))
                if (!$util.isString(message.body))
                    return "body: string expected";
            return null;
        };
        /**
         * Creates a Message message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.Message
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.Message} Message
         */
        Message.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.Message)
                return object;
            let message = new $root.messagespb.Message();
            if (object.id != null)
                message.id = String(object.id);
            if (object.conversationId != null)
                message.conversationId = String(object.conversationId);
            if (object.senderId != null)
                message.senderId = String(object.senderId);
            if (object.recipientId != null)
                message.recipientId = String(object.recipientId);
            if (object.itemId != null)
                message.itemId = String(object.itemId);
            if (object.body != null)
                message.body = String(object.body);
            return message;
        };
        /**
         * Creates a plain object from a Message message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.Message
         * @static
         * @param {messagespb.Message} message Message
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        Message.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.id = "";
                object.conversationId = "";
                object.senderId = "";
                object.recipientId = "";
                object.itemId = "";
                object.body = "";
            }
            if (message.id != null && message.hasOwnProperty("id"))
                object.id = message.id;
            if (message.conversationId != null && message.hasOwnProperty("conversationId"))
                object.conversationId = message.conversationId;
            if (message.senderId != null && message.hasOwnProperty("senderId"))
                object.senderId = message.senderId;
            if (message.recipientId != null && message.hasOwnProperty("recipientId"))
                object.recipientId = message.recipientId;
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                object.itemId = message.itemId;
            if (message.body != null && message.hasOwnProperty("body"))
                object.body = message.body;
            return object;
        };
        /**
         * Converts this Message to JSON.
         * @function toJSON
         * @memberof messagespb.Message
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        Message.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for Message
         * @function getTypeUrl
         * @memberof messagespb.Message
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        Message.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.Message";
        };
        return Message;
    })();
    /**
     * ConversationStatus enum.
     * @name messagespb.ConversationStatus
     * @enum {number}
     * @property {number} PENDING=0 PENDING value
     * @property {number} ACTIVE=1 ACTIVE value
     * @property {number} ARCHIVED=2 ARCHIVED value
     */
    messagespb.ConversationStatus = (function() {
        const valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "PENDING"] = 0;
        values[valuesById[1] = "ACTIVE"] = 1;
        values[valuesById[2] = "ARCHIVED"] = 2;
        return values;
    })();
    /**
     * MessageStatus enum.
     * @name messagespb.MessageStatus
     * @enum {number}
     * @property {number} SENT=0 SENT value
     * @property {number} EDITED=1 EDITED value
     * @property {number} DELETED=2 DELETED value
     */
    messagespb.MessageStatus = (function() {
        const valuesById = {}, values = Object.create(valuesById);
        values[valuesById[0] = "SENT"] = 0;
        values[valuesById[1] = "EDITED"] = 1;
        values[valuesById[2] = "DELETED"] = 2;
        return values;
    })();
    messagespb.StartConversationRequest = (function() {
        /**
         * Properties of a StartConversationRequest.
         * @memberof messagespb
         * @interface IStartConversationRequest
         * @property {string|null} [senderId] StartConversationRequest senderId
         * @property {string|null} [recipientId] StartConversationRequest recipientId
         * @property {string|null} [itemId] StartConversationRequest itemId
         */
        /**
         * Constructs a new StartConversationRequest.
         * @memberof messagespb
         * @classdesc Represents a StartConversationRequest.
         * @implements IStartConversationRequest
         * @constructor
         * @param {messagespb.IStartConversationRequest=} [properties] Properties to set
         */
        function StartConversationRequest(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * StartConversationRequest senderId.
         * @member {string} senderId
         * @memberof messagespb.StartConversationRequest
         * @instance
         */
        StartConversationRequest.prototype.senderId = "";
        /**
         * StartConversationRequest recipientId.
         * @member {string} recipientId
         * @memberof messagespb.StartConversationRequest
         * @instance
         */
        StartConversationRequest.prototype.recipientId = "";
        /**
         * StartConversationRequest itemId.
         * @member {string} itemId
         * @memberof messagespb.StartConversationRequest
         * @instance
         */
        StartConversationRequest.prototype.itemId = "";
        /**
         * Creates a new StartConversationRequest instance using the specified properties.
         * @function create
         * @memberof messagespb.StartConversationRequest
         * @static
         * @param {messagespb.IStartConversationRequest=} [properties] Properties to set
         * @returns {messagespb.StartConversationRequest} StartConversationRequest instance
         */
        StartConversationRequest.create = function create(properties) {
            return new StartConversationRequest(properties);
        };
        /**
         * Encodes the specified StartConversationRequest message. Does not implicitly {@link messagespb.StartConversationRequest.verify|verify} messages.
         * @function encode
         * @memberof messagespb.StartConversationRequest
         * @static
         * @param {messagespb.IStartConversationRequest} message StartConversationRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        StartConversationRequest.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.senderId != null && Object.hasOwnProperty.call(message, "senderId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.senderId);
            if (message.recipientId != null && Object.hasOwnProperty.call(message, "recipientId"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.recipientId);
            if (message.itemId != null && Object.hasOwnProperty.call(message, "itemId"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.itemId);
            return writer;
        };
        /**
         * Encodes the specified StartConversationRequest message, length delimited. Does not implicitly {@link messagespb.StartConversationRequest.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.StartConversationRequest
         * @static
         * @param {messagespb.IStartConversationRequest} message StartConversationRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        StartConversationRequest.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a StartConversationRequest message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.StartConversationRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.StartConversationRequest} StartConversationRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        StartConversationRequest.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.StartConversationRequest();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.senderId = reader.string();
                        break;
                    }
                case 2: {
                        message.recipientId = reader.string();
                        break;
                    }
                case 3: {
                        message.itemId = reader.string();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a StartConversationRequest message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.StartConversationRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.StartConversationRequest} StartConversationRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        StartConversationRequest.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a StartConversationRequest message.
         * @function verify
         * @memberof messagespb.StartConversationRequest
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        StartConversationRequest.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.senderId != null && message.hasOwnProperty("senderId"))
                if (!$util.isString(message.senderId))
                    return "senderId: string expected";
            if (message.recipientId != null && message.hasOwnProperty("recipientId"))
                if (!$util.isString(message.recipientId))
                    return "recipientId: string expected";
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                if (!$util.isString(message.itemId))
                    return "itemId: string expected";
            return null;
        };
        /**
         * Creates a StartConversationRequest message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.StartConversationRequest
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.StartConversationRequest} StartConversationRequest
         */
        StartConversationRequest.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.StartConversationRequest)
                return object;
            let message = new $root.messagespb.StartConversationRequest();
            if (object.senderId != null)
                message.senderId = String(object.senderId);
            if (object.recipientId != null)
                message.recipientId = String(object.recipientId);
            if (object.itemId != null)
                message.itemId = String(object.itemId);
            return message;
        };
        /**
         * Creates a plain object from a StartConversationRequest message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.StartConversationRequest
         * @static
         * @param {messagespb.StartConversationRequest} message StartConversationRequest
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        StartConversationRequest.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.senderId = "";
                object.recipientId = "";
                object.itemId = "";
            }
            if (message.senderId != null && message.hasOwnProperty("senderId"))
                object.senderId = message.senderId;
            if (message.recipientId != null && message.hasOwnProperty("recipientId"))
                object.recipientId = message.recipientId;
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                object.itemId = message.itemId;
            return object;
        };
        /**
         * Converts this StartConversationRequest to JSON.
         * @function toJSON
         * @memberof messagespb.StartConversationRequest
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        StartConversationRequest.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for StartConversationRequest
         * @function getTypeUrl
         * @memberof messagespb.StartConversationRequest
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        StartConversationRequest.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.StartConversationRequest";
        };
        return StartConversationRequest;
    })();
    messagespb.StartConversationResponse = (function() {
        /**
         * Properties of a StartConversationResponse.
         * @memberof messagespb
         * @interface IStartConversationResponse
         * @property {string|null} [id] StartConversationResponse id
         */
        /**
         * Constructs a new StartConversationResponse.
         * @memberof messagespb
         * @classdesc Represents a StartConversationResponse.
         * @implements IStartConversationResponse
         * @constructor
         * @param {messagespb.IStartConversationResponse=} [properties] Properties to set
         */
        function StartConversationResponse(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * StartConversationResponse id.
         * @member {string} id
         * @memberof messagespb.StartConversationResponse
         * @instance
         */
        StartConversationResponse.prototype.id = "";
        /**
         * Creates a new StartConversationResponse instance using the specified properties.
         * @function create
         * @memberof messagespb.StartConversationResponse
         * @static
         * @param {messagespb.IStartConversationResponse=} [properties] Properties to set
         * @returns {messagespb.StartConversationResponse} StartConversationResponse instance
         */
        StartConversationResponse.create = function create(properties) {
            return new StartConversationResponse(properties);
        };
        /**
         * Encodes the specified StartConversationResponse message. Does not implicitly {@link messagespb.StartConversationResponse.verify|verify} messages.
         * @function encode
         * @memberof messagespb.StartConversationResponse
         * @static
         * @param {messagespb.IStartConversationResponse} message StartConversationResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        StartConversationResponse.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
            return writer;
        };
        /**
         * Encodes the specified StartConversationResponse message, length delimited. Does not implicitly {@link messagespb.StartConversationResponse.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.StartConversationResponse
         * @static
         * @param {messagespb.IStartConversationResponse} message StartConversationResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        StartConversationResponse.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a StartConversationResponse message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.StartConversationResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.StartConversationResponse} StartConversationResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        StartConversationResponse.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.StartConversationResponse();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.id = reader.string();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a StartConversationResponse message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.StartConversationResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.StartConversationResponse} StartConversationResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        StartConversationResponse.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a StartConversationResponse message.
         * @function verify
         * @memberof messagespb.StartConversationResponse
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        StartConversationResponse.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.id != null && message.hasOwnProperty("id"))
                if (!$util.isString(message.id))
                    return "id: string expected";
            return null;
        };
        /**
         * Creates a StartConversationResponse message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.StartConversationResponse
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.StartConversationResponse} StartConversationResponse
         */
        StartConversationResponse.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.StartConversationResponse)
                return object;
            let message = new $root.messagespb.StartConversationResponse();
            if (object.id != null)
                message.id = String(object.id);
            return message;
        };
        /**
         * Creates a plain object from a StartConversationResponse message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.StartConversationResponse
         * @static
         * @param {messagespb.StartConversationResponse} message StartConversationResponse
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        StartConversationResponse.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults)
                object.id = "";
            if (message.id != null && message.hasOwnProperty("id"))
                object.id = message.id;
            return object;
        };
        /**
         * Converts this StartConversationResponse to JSON.
         * @function toJSON
         * @memberof messagespb.StartConversationResponse
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        StartConversationResponse.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for StartConversationResponse
         * @function getTypeUrl
         * @memberof messagespb.StartConversationResponse
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        StartConversationResponse.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.StartConversationResponse";
        };
        return StartConversationResponse;
    })();
    messagespb.RestoreConversationRequest = (function() {
        /**
         * Properties of a RestoreConversationRequest.
         * @memberof messagespb
         * @interface IRestoreConversationRequest
         * @property {string|null} [id] RestoreConversationRequest id
         */
        /**
         * Constructs a new RestoreConversationRequest.
         * @memberof messagespb
         * @classdesc Represents a RestoreConversationRequest.
         * @implements IRestoreConversationRequest
         * @constructor
         * @param {messagespb.IRestoreConversationRequest=} [properties] Properties to set
         */
        function RestoreConversationRequest(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * RestoreConversationRequest id.
         * @member {string} id
         * @memberof messagespb.RestoreConversationRequest
         * @instance
         */
        RestoreConversationRequest.prototype.id = "";
        /**
         * Creates a new RestoreConversationRequest instance using the specified properties.
         * @function create
         * @memberof messagespb.RestoreConversationRequest
         * @static
         * @param {messagespb.IRestoreConversationRequest=} [properties] Properties to set
         * @returns {messagespb.RestoreConversationRequest} RestoreConversationRequest instance
         */
        RestoreConversationRequest.create = function create(properties) {
            return new RestoreConversationRequest(properties);
        };
        /**
         * Encodes the specified RestoreConversationRequest message. Does not implicitly {@link messagespb.RestoreConversationRequest.verify|verify} messages.
         * @function encode
         * @memberof messagespb.RestoreConversationRequest
         * @static
         * @param {messagespb.IRestoreConversationRequest} message RestoreConversationRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RestoreConversationRequest.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
            return writer;
        };
        /**
         * Encodes the specified RestoreConversationRequest message, length delimited. Does not implicitly {@link messagespb.RestoreConversationRequest.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.RestoreConversationRequest
         * @static
         * @param {messagespb.IRestoreConversationRequest} message RestoreConversationRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RestoreConversationRequest.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a RestoreConversationRequest message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.RestoreConversationRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.RestoreConversationRequest} RestoreConversationRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RestoreConversationRequest.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.RestoreConversationRequest();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.id = reader.string();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a RestoreConversationRequest message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.RestoreConversationRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.RestoreConversationRequest} RestoreConversationRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RestoreConversationRequest.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a RestoreConversationRequest message.
         * @function verify
         * @memberof messagespb.RestoreConversationRequest
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        RestoreConversationRequest.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.id != null && message.hasOwnProperty("id"))
                if (!$util.isString(message.id))
                    return "id: string expected";
            return null;
        };
        /**
         * Creates a RestoreConversationRequest message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.RestoreConversationRequest
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.RestoreConversationRequest} RestoreConversationRequest
         */
        RestoreConversationRequest.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.RestoreConversationRequest)
                return object;
            let message = new $root.messagespb.RestoreConversationRequest();
            if (object.id != null)
                message.id = String(object.id);
            return message;
        };
        /**
         * Creates a plain object from a RestoreConversationRequest message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.RestoreConversationRequest
         * @static
         * @param {messagespb.RestoreConversationRequest} message RestoreConversationRequest
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        RestoreConversationRequest.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults)
                object.id = "";
            if (message.id != null && message.hasOwnProperty("id"))
                object.id = message.id;
            return object;
        };
        /**
         * Converts this RestoreConversationRequest to JSON.
         * @function toJSON
         * @memberof messagespb.RestoreConversationRequest
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        RestoreConversationRequest.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for RestoreConversationRequest
         * @function getTypeUrl
         * @memberof messagespb.RestoreConversationRequest
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        RestoreConversationRequest.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.RestoreConversationRequest";
        };
        return RestoreConversationRequest;
    })();
    messagespb.RestoreConversationResponse = (function() {
        /**
         * Properties of a RestoreConversationResponse.
         * @memberof messagespb
         * @interface IRestoreConversationResponse
         * @property {string|null} [id] RestoreConversationResponse id
         * @property {string|null} [conversationStatus] RestoreConversationResponse conversationStatus
         */
        /**
         * Constructs a new RestoreConversationResponse.
         * @memberof messagespb
         * @classdesc Represents a RestoreConversationResponse.
         * @implements IRestoreConversationResponse
         * @constructor
         * @param {messagespb.IRestoreConversationResponse=} [properties] Properties to set
         */
        function RestoreConversationResponse(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * RestoreConversationResponse id.
         * @member {string} id
         * @memberof messagespb.RestoreConversationResponse
         * @instance
         */
        RestoreConversationResponse.prototype.id = "";
        /**
         * RestoreConversationResponse conversationStatus.
         * @member {string} conversationStatus
         * @memberof messagespb.RestoreConversationResponse
         * @instance
         */
        RestoreConversationResponse.prototype.conversationStatus = "";
        /**
         * Creates a new RestoreConversationResponse instance using the specified properties.
         * @function create
         * @memberof messagespb.RestoreConversationResponse
         * @static
         * @param {messagespb.IRestoreConversationResponse=} [properties] Properties to set
         * @returns {messagespb.RestoreConversationResponse} RestoreConversationResponse instance
         */
        RestoreConversationResponse.create = function create(properties) {
            return new RestoreConversationResponse(properties);
        };
        /**
         * Encodes the specified RestoreConversationResponse message. Does not implicitly {@link messagespb.RestoreConversationResponse.verify|verify} messages.
         * @function encode
         * @memberof messagespb.RestoreConversationResponse
         * @static
         * @param {messagespb.IRestoreConversationResponse} message RestoreConversationResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RestoreConversationResponse.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
            if (message.conversationStatus != null && Object.hasOwnProperty.call(message, "conversationStatus"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.conversationStatus);
            return writer;
        };
        /**
         * Encodes the specified RestoreConversationResponse message, length delimited. Does not implicitly {@link messagespb.RestoreConversationResponse.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.RestoreConversationResponse
         * @static
         * @param {messagespb.IRestoreConversationResponse} message RestoreConversationResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RestoreConversationResponse.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a RestoreConversationResponse message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.RestoreConversationResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.RestoreConversationResponse} RestoreConversationResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RestoreConversationResponse.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.RestoreConversationResponse();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.id = reader.string();
                        break;
                    }
                case 2: {
                        message.conversationStatus = reader.string();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a RestoreConversationResponse message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.RestoreConversationResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.RestoreConversationResponse} RestoreConversationResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RestoreConversationResponse.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a RestoreConversationResponse message.
         * @function verify
         * @memberof messagespb.RestoreConversationResponse
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        RestoreConversationResponse.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.id != null && message.hasOwnProperty("id"))
                if (!$util.isString(message.id))
                    return "id: string expected";
            if (message.conversationStatus != null && message.hasOwnProperty("conversationStatus"))
                if (!$util.isString(message.conversationStatus))
                    return "conversationStatus: string expected";
            return null;
        };
        /**
         * Creates a RestoreConversationResponse message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.RestoreConversationResponse
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.RestoreConversationResponse} RestoreConversationResponse
         */
        RestoreConversationResponse.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.RestoreConversationResponse)
                return object;
            let message = new $root.messagespb.RestoreConversationResponse();
            if (object.id != null)
                message.id = String(object.id);
            if (object.conversationStatus != null)
                message.conversationStatus = String(object.conversationStatus);
            return message;
        };
        /**
         * Creates a plain object from a RestoreConversationResponse message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.RestoreConversationResponse
         * @static
         * @param {messagespb.RestoreConversationResponse} message RestoreConversationResponse
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        RestoreConversationResponse.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.id = "";
                object.conversationStatus = "";
            }
            if (message.id != null && message.hasOwnProperty("id"))
                object.id = message.id;
            if (message.conversationStatus != null && message.hasOwnProperty("conversationStatus"))
                object.conversationStatus = message.conversationStatus;
            return object;
        };
        /**
         * Converts this RestoreConversationResponse to JSON.
         * @function toJSON
         * @memberof messagespb.RestoreConversationResponse
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        RestoreConversationResponse.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for RestoreConversationResponse
         * @function getTypeUrl
         * @memberof messagespb.RestoreConversationResponse
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        RestoreConversationResponse.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.RestoreConversationResponse";
        };
        return RestoreConversationResponse;
    })();
    messagespb.ArchiveConversationRequest = (function() {
        /**
         * Properties of an ArchiveConversationRequest.
         * @memberof messagespb
         * @interface IArchiveConversationRequest
         * @property {string|null} [id] ArchiveConversationRequest id
         */
        /**
         * Constructs a new ArchiveConversationRequest.
         * @memberof messagespb
         * @classdesc Represents an ArchiveConversationRequest.
         * @implements IArchiveConversationRequest
         * @constructor
         * @param {messagespb.IArchiveConversationRequest=} [properties] Properties to set
         */
        function ArchiveConversationRequest(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * ArchiveConversationRequest id.
         * @member {string} id
         * @memberof messagespb.ArchiveConversationRequest
         * @instance
         */
        ArchiveConversationRequest.prototype.id = "";
        /**
         * Creates a new ArchiveConversationRequest instance using the specified properties.
         * @function create
         * @memberof messagespb.ArchiveConversationRequest
         * @static
         * @param {messagespb.IArchiveConversationRequest=} [properties] Properties to set
         * @returns {messagespb.ArchiveConversationRequest} ArchiveConversationRequest instance
         */
        ArchiveConversationRequest.create = function create(properties) {
            return new ArchiveConversationRequest(properties);
        };
        /**
         * Encodes the specified ArchiveConversationRequest message. Does not implicitly {@link messagespb.ArchiveConversationRequest.verify|verify} messages.
         * @function encode
         * @memberof messagespb.ArchiveConversationRequest
         * @static
         * @param {messagespb.IArchiveConversationRequest} message ArchiveConversationRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ArchiveConversationRequest.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
            return writer;
        };
        /**
         * Encodes the specified ArchiveConversationRequest message, length delimited. Does not implicitly {@link messagespb.ArchiveConversationRequest.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.ArchiveConversationRequest
         * @static
         * @param {messagespb.IArchiveConversationRequest} message ArchiveConversationRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ArchiveConversationRequest.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes an ArchiveConversationRequest message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.ArchiveConversationRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.ArchiveConversationRequest} ArchiveConversationRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ArchiveConversationRequest.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.ArchiveConversationRequest();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.id = reader.string();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes an ArchiveConversationRequest message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.ArchiveConversationRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.ArchiveConversationRequest} ArchiveConversationRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ArchiveConversationRequest.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies an ArchiveConversationRequest message.
         * @function verify
         * @memberof messagespb.ArchiveConversationRequest
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        ArchiveConversationRequest.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.id != null && message.hasOwnProperty("id"))
                if (!$util.isString(message.id))
                    return "id: string expected";
            return null;
        };
        /**
         * Creates an ArchiveConversationRequest message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.ArchiveConversationRequest
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.ArchiveConversationRequest} ArchiveConversationRequest
         */
        ArchiveConversationRequest.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.ArchiveConversationRequest)
                return object;
            let message = new $root.messagespb.ArchiveConversationRequest();
            if (object.id != null)
                message.id = String(object.id);
            return message;
        };
        /**
         * Creates a plain object from an ArchiveConversationRequest message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.ArchiveConversationRequest
         * @static
         * @param {messagespb.ArchiveConversationRequest} message ArchiveConversationRequest
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        ArchiveConversationRequest.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults)
                object.id = "";
            if (message.id != null && message.hasOwnProperty("id"))
                object.id = message.id;
            return object;
        };
        /**
         * Converts this ArchiveConversationRequest to JSON.
         * @function toJSON
         * @memberof messagespb.ArchiveConversationRequest
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        ArchiveConversationRequest.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for ArchiveConversationRequest
         * @function getTypeUrl
         * @memberof messagespb.ArchiveConversationRequest
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        ArchiveConversationRequest.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.ArchiveConversationRequest";
        };
        return ArchiveConversationRequest;
    })();
    messagespb.ArchiveConversationResponse = (function() {
        /**
         * Properties of an ArchiveConversationResponse.
         * @memberof messagespb
         * @interface IArchiveConversationResponse
         * @property {string|null} [id] ArchiveConversationResponse id
         * @property {string|null} [conversationStatus] ArchiveConversationResponse conversationStatus
         */
        /**
         * Constructs a new ArchiveConversationResponse.
         * @memberof messagespb
         * @classdesc Represents an ArchiveConversationResponse.
         * @implements IArchiveConversationResponse
         * @constructor
         * @param {messagespb.IArchiveConversationResponse=} [properties] Properties to set
         */
        function ArchiveConversationResponse(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * ArchiveConversationResponse id.
         * @member {string} id
         * @memberof messagespb.ArchiveConversationResponse
         * @instance
         */
        ArchiveConversationResponse.prototype.id = "";
        /**
         * ArchiveConversationResponse conversationStatus.
         * @member {string} conversationStatus
         * @memberof messagespb.ArchiveConversationResponse
         * @instance
         */
        ArchiveConversationResponse.prototype.conversationStatus = "";
        /**
         * Creates a new ArchiveConversationResponse instance using the specified properties.
         * @function create
         * @memberof messagespb.ArchiveConversationResponse
         * @static
         * @param {messagespb.IArchiveConversationResponse=} [properties] Properties to set
         * @returns {messagespb.ArchiveConversationResponse} ArchiveConversationResponse instance
         */
        ArchiveConversationResponse.create = function create(properties) {
            return new ArchiveConversationResponse(properties);
        };
        /**
         * Encodes the specified ArchiveConversationResponse message. Does not implicitly {@link messagespb.ArchiveConversationResponse.verify|verify} messages.
         * @function encode
         * @memberof messagespb.ArchiveConversationResponse
         * @static
         * @param {messagespb.IArchiveConversationResponse} message ArchiveConversationResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ArchiveConversationResponse.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
            if (message.conversationStatus != null && Object.hasOwnProperty.call(message, "conversationStatus"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.conversationStatus);
            return writer;
        };
        /**
         * Encodes the specified ArchiveConversationResponse message, length delimited. Does not implicitly {@link messagespb.ArchiveConversationResponse.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.ArchiveConversationResponse
         * @static
         * @param {messagespb.IArchiveConversationResponse} message ArchiveConversationResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ArchiveConversationResponse.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes an ArchiveConversationResponse message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.ArchiveConversationResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.ArchiveConversationResponse} ArchiveConversationResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ArchiveConversationResponse.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.ArchiveConversationResponse();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.id = reader.string();
                        break;
                    }
                case 2: {
                        message.conversationStatus = reader.string();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes an ArchiveConversationResponse message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.ArchiveConversationResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.ArchiveConversationResponse} ArchiveConversationResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ArchiveConversationResponse.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies an ArchiveConversationResponse message.
         * @function verify
         * @memberof messagespb.ArchiveConversationResponse
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        ArchiveConversationResponse.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.id != null && message.hasOwnProperty("id"))
                if (!$util.isString(message.id))
                    return "id: string expected";
            if (message.conversationStatus != null && message.hasOwnProperty("conversationStatus"))
                if (!$util.isString(message.conversationStatus))
                    return "conversationStatus: string expected";
            return null;
        };
        /**
         * Creates an ArchiveConversationResponse message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.ArchiveConversationResponse
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.ArchiveConversationResponse} ArchiveConversationResponse
         */
        ArchiveConversationResponse.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.ArchiveConversationResponse)
                return object;
            let message = new $root.messagespb.ArchiveConversationResponse();
            if (object.id != null)
                message.id = String(object.id);
            if (object.conversationStatus != null)
                message.conversationStatus = String(object.conversationStatus);
            return message;
        };
        /**
         * Creates a plain object from an ArchiveConversationResponse message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.ArchiveConversationResponse
         * @static
         * @param {messagespb.ArchiveConversationResponse} message ArchiveConversationResponse
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        ArchiveConversationResponse.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.id = "";
                object.conversationStatus = "";
            }
            if (message.id != null && message.hasOwnProperty("id"))
                object.id = message.id;
            if (message.conversationStatus != null && message.hasOwnProperty("conversationStatus"))
                object.conversationStatus = message.conversationStatus;
            return object;
        };
        /**
         * Converts this ArchiveConversationResponse to JSON.
         * @function toJSON
         * @memberof messagespb.ArchiveConversationResponse
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        ArchiveConversationResponse.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for ArchiveConversationResponse
         * @function getTypeUrl
         * @memberof messagespb.ArchiveConversationResponse
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        ArchiveConversationResponse.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.ArchiveConversationResponse";
        };
        return ArchiveConversationResponse;
    })();
    messagespb.GetConversationRequest = (function() {
        /**
         * Properties of a GetConversationRequest.
         * @memberof messagespb
         * @interface IGetConversationRequest
         * @property {string|null} [id] GetConversationRequest id
         */
        /**
         * Constructs a new GetConversationRequest.
         * @memberof messagespb
         * @classdesc Represents a GetConversationRequest.
         * @implements IGetConversationRequest
         * @constructor
         * @param {messagespb.IGetConversationRequest=} [properties] Properties to set
         */
        function GetConversationRequest(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * GetConversationRequest id.
         * @member {string} id
         * @memberof messagespb.GetConversationRequest
         * @instance
         */
        GetConversationRequest.prototype.id = "";
        /**
         * Creates a new GetConversationRequest instance using the specified properties.
         * @function create
         * @memberof messagespb.GetConversationRequest
         * @static
         * @param {messagespb.IGetConversationRequest=} [properties] Properties to set
         * @returns {messagespb.GetConversationRequest} GetConversationRequest instance
         */
        GetConversationRequest.create = function create(properties) {
            return new GetConversationRequest(properties);
        };
        /**
         * Encodes the specified GetConversationRequest message. Does not implicitly {@link messagespb.GetConversationRequest.verify|verify} messages.
         * @function encode
         * @memberof messagespb.GetConversationRequest
         * @static
         * @param {messagespb.IGetConversationRequest} message GetConversationRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetConversationRequest.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
            return writer;
        };
        /**
         * Encodes the specified GetConversationRequest message, length delimited. Does not implicitly {@link messagespb.GetConversationRequest.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.GetConversationRequest
         * @static
         * @param {messagespb.IGetConversationRequest} message GetConversationRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetConversationRequest.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a GetConversationRequest message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.GetConversationRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.GetConversationRequest} GetConversationRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetConversationRequest.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.GetConversationRequest();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.id = reader.string();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a GetConversationRequest message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.GetConversationRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.GetConversationRequest} GetConversationRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetConversationRequest.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a GetConversationRequest message.
         * @function verify
         * @memberof messagespb.GetConversationRequest
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        GetConversationRequest.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.id != null && message.hasOwnProperty("id"))
                if (!$util.isString(message.id))
                    return "id: string expected";
            return null;
        };
        /**
         * Creates a GetConversationRequest message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.GetConversationRequest
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.GetConversationRequest} GetConversationRequest
         */
        GetConversationRequest.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.GetConversationRequest)
                return object;
            let message = new $root.messagespb.GetConversationRequest();
            if (object.id != null)
                message.id = String(object.id);
            return message;
        };
        /**
         * Creates a plain object from a GetConversationRequest message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.GetConversationRequest
         * @static
         * @param {messagespb.GetConversationRequest} message GetConversationRequest
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        GetConversationRequest.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults)
                object.id = "";
            if (message.id != null && message.hasOwnProperty("id"))
                object.id = message.id;
            return object;
        };
        /**
         * Converts this GetConversationRequest to JSON.
         * @function toJSON
         * @memberof messagespb.GetConversationRequest
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        GetConversationRequest.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for GetConversationRequest
         * @function getTypeUrl
         * @memberof messagespb.GetConversationRequest
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        GetConversationRequest.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.GetConversationRequest";
        };
        return GetConversationRequest;
    })();
    messagespb.GetConversationResponse = (function() {
        /**
         * Properties of a GetConversationResponse.
         * @memberof messagespb
         * @interface IGetConversationResponse
         * @property {messagespb.IConversation|null} [conversation] GetConversationResponse conversation
         */
        /**
         * Constructs a new GetConversationResponse.
         * @memberof messagespb
         * @classdesc Represents a GetConversationResponse.
         * @implements IGetConversationResponse
         * @constructor
         * @param {messagespb.IGetConversationResponse=} [properties] Properties to set
         */
        function GetConversationResponse(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * GetConversationResponse conversation.
         * @member {messagespb.IConversation|null|undefined} conversation
         * @memberof messagespb.GetConversationResponse
         * @instance
         */
        GetConversationResponse.prototype.conversation = null;
        /**
         * Creates a new GetConversationResponse instance using the specified properties.
         * @function create
         * @memberof messagespb.GetConversationResponse
         * @static
         * @param {messagespb.IGetConversationResponse=} [properties] Properties to set
         * @returns {messagespb.GetConversationResponse} GetConversationResponse instance
         */
        GetConversationResponse.create = function create(properties) {
            return new GetConversationResponse(properties);
        };
        /**
         * Encodes the specified GetConversationResponse message. Does not implicitly {@link messagespb.GetConversationResponse.verify|verify} messages.
         * @function encode
         * @memberof messagespb.GetConversationResponse
         * @static
         * @param {messagespb.IGetConversationResponse} message GetConversationResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetConversationResponse.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.conversation != null && Object.hasOwnProperty.call(message, "conversation"))
                $root.messagespb.Conversation.encode(message.conversation, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            return writer;
        };
        /**
         * Encodes the specified GetConversationResponse message, length delimited. Does not implicitly {@link messagespb.GetConversationResponse.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.GetConversationResponse
         * @static
         * @param {messagespb.IGetConversationResponse} message GetConversationResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetConversationResponse.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a GetConversationResponse message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.GetConversationResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.GetConversationResponse} GetConversationResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetConversationResponse.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.GetConversationResponse();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.conversation = $root.messagespb.Conversation.decode(reader, reader.uint32());
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a GetConversationResponse message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.GetConversationResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.GetConversationResponse} GetConversationResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetConversationResponse.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a GetConversationResponse message.
         * @function verify
         * @memberof messagespb.GetConversationResponse
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        GetConversationResponse.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.conversation != null && message.hasOwnProperty("conversation")) {
                let error = $root.messagespb.Conversation.verify(message.conversation);
                if (error)
                    return "conversation." + error;
            }
            return null;
        };
        /**
         * Creates a GetConversationResponse message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.GetConversationResponse
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.GetConversationResponse} GetConversationResponse
         */
        GetConversationResponse.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.GetConversationResponse)
                return object;
            let message = new $root.messagespb.GetConversationResponse();
            if (object.conversation != null) {
                if (typeof object.conversation !== "object")
                    throw TypeError(".messagespb.GetConversationResponse.conversation: object expected");
                message.conversation = $root.messagespb.Conversation.fromObject(object.conversation);
            }
            return message;
        };
        /**
         * Creates a plain object from a GetConversationResponse message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.GetConversationResponse
         * @static
         * @param {messagespb.GetConversationResponse} message GetConversationResponse
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        GetConversationResponse.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults)
                object.conversation = null;
            if (message.conversation != null && message.hasOwnProperty("conversation"))
                object.conversation = $root.messagespb.Conversation.toObject(message.conversation, options);
            return object;
        };
        /**
         * Converts this GetConversationResponse to JSON.
         * @function toJSON
         * @memberof messagespb.GetConversationResponse
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        GetConversationResponse.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for GetConversationResponse
         * @function getTypeUrl
         * @memberof messagespb.GetConversationResponse
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        GetConversationResponse.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.GetConversationResponse";
        };
        return GetConversationResponse;
    })();
    messagespb.GetConversationsRequest = (function() {
        /**
         * Properties of a GetConversationsRequest.
         * @memberof messagespb
         * @interface IGetConversationsRequest
         * @property {string|null} [userId] GetConversationsRequest userId
         * @property {number|null} [page] GetConversationsRequest page
         * @property {number|null} [limit] GetConversationsRequest limit
         */
        /**
         * Constructs a new GetConversationsRequest.
         * @memberof messagespb
         * @classdesc Represents a GetConversationsRequest.
         * @implements IGetConversationsRequest
         * @constructor
         * @param {messagespb.IGetConversationsRequest=} [properties] Properties to set
         */
        function GetConversationsRequest(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * GetConversationsRequest userId.
         * @member {string} userId
         * @memberof messagespb.GetConversationsRequest
         * @instance
         */
        GetConversationsRequest.prototype.userId = "";
        /**
         * GetConversationsRequest page.
         * @member {number} page
         * @memberof messagespb.GetConversationsRequest
         * @instance
         */
        GetConversationsRequest.prototype.page = 0;
        /**
         * GetConversationsRequest limit.
         * @member {number} limit
         * @memberof messagespb.GetConversationsRequest
         * @instance
         */
        GetConversationsRequest.prototype.limit = 0;
        /**
         * Creates a new GetConversationsRequest instance using the specified properties.
         * @function create
         * @memberof messagespb.GetConversationsRequest
         * @static
         * @param {messagespb.IGetConversationsRequest=} [properties] Properties to set
         * @returns {messagespb.GetConversationsRequest} GetConversationsRequest instance
         */
        GetConversationsRequest.create = function create(properties) {
            return new GetConversationsRequest(properties);
        };
        /**
         * Encodes the specified GetConversationsRequest message. Does not implicitly {@link messagespb.GetConversationsRequest.verify|verify} messages.
         * @function encode
         * @memberof messagespb.GetConversationsRequest
         * @static
         * @param {messagespb.IGetConversationsRequest} message GetConversationsRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetConversationsRequest.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.userId != null && Object.hasOwnProperty.call(message, "userId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.userId);
            if (message.page != null && Object.hasOwnProperty.call(message, "page"))
                writer.uint32(/* id 2, wireType 0 =*/16).int32(message.page);
            if (message.limit != null && Object.hasOwnProperty.call(message, "limit"))
                writer.uint32(/* id 3, wireType 0 =*/24).int32(message.limit);
            return writer;
        };
        /**
         * Encodes the specified GetConversationsRequest message, length delimited. Does not implicitly {@link messagespb.GetConversationsRequest.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.GetConversationsRequest
         * @static
         * @param {messagespb.IGetConversationsRequest} message GetConversationsRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetConversationsRequest.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a GetConversationsRequest message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.GetConversationsRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.GetConversationsRequest} GetConversationsRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetConversationsRequest.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.GetConversationsRequest();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.userId = reader.string();
                        break;
                    }
                case 2: {
                        message.page = reader.int32();
                        break;
                    }
                case 3: {
                        message.limit = reader.int32();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a GetConversationsRequest message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.GetConversationsRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.GetConversationsRequest} GetConversationsRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetConversationsRequest.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a GetConversationsRequest message.
         * @function verify
         * @memberof messagespb.GetConversationsRequest
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        GetConversationsRequest.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.userId != null && message.hasOwnProperty("userId"))
                if (!$util.isString(message.userId))
                    return "userId: string expected";
            if (message.page != null && message.hasOwnProperty("page"))
                if (!$util.isInteger(message.page))
                    return "page: integer expected";
            if (message.limit != null && message.hasOwnProperty("limit"))
                if (!$util.isInteger(message.limit))
                    return "limit: integer expected";
            return null;
        };
        /**
         * Creates a GetConversationsRequest message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.GetConversationsRequest
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.GetConversationsRequest} GetConversationsRequest
         */
        GetConversationsRequest.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.GetConversationsRequest)
                return object;
            let message = new $root.messagespb.GetConversationsRequest();
            if (object.userId != null)
                message.userId = String(object.userId);
            if (object.page != null)
                message.page = object.page | 0;
            if (object.limit != null)
                message.limit = object.limit | 0;
            return message;
        };
        /**
         * Creates a plain object from a GetConversationsRequest message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.GetConversationsRequest
         * @static
         * @param {messagespb.GetConversationsRequest} message GetConversationsRequest
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        GetConversationsRequest.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.userId = "";
                object.page = 0;
                object.limit = 0;
            }
            if (message.userId != null && message.hasOwnProperty("userId"))
                object.userId = message.userId;
            if (message.page != null && message.hasOwnProperty("page"))
                object.page = message.page;
            if (message.limit != null && message.hasOwnProperty("limit"))
                object.limit = message.limit;
            return object;
        };
        /**
         * Converts this GetConversationsRequest to JSON.
         * @function toJSON
         * @memberof messagespb.GetConversationsRequest
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        GetConversationsRequest.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for GetConversationsRequest
         * @function getTypeUrl
         * @memberof messagespb.GetConversationsRequest
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        GetConversationsRequest.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.GetConversationsRequest";
        };
        return GetConversationsRequest;
    })();
    messagespb.GetConversationsResponse = (function() {
        /**
         * Properties of a GetConversationsResponse.
         * @memberof messagespb
         * @interface IGetConversationsResponse
         * @property {Array.<messagespb.IConversation>|null} [conversations] GetConversationsResponse conversations
         * @property {number|null} [total] GetConversationsResponse total
         * @property {number|null} [page] GetConversationsResponse page
         * @property {number|null} [limit] GetConversationsResponse limit
         */
        /**
         * Constructs a new GetConversationsResponse.
         * @memberof messagespb
         * @classdesc Represents a GetConversationsResponse.
         * @implements IGetConversationsResponse
         * @constructor
         * @param {messagespb.IGetConversationsResponse=} [properties] Properties to set
         */
        function GetConversationsResponse(properties) {
            this.conversations = [];
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * GetConversationsResponse conversations.
         * @member {Array.<messagespb.IConversation>} conversations
         * @memberof messagespb.GetConversationsResponse
         * @instance
         */
        GetConversationsResponse.prototype.conversations = $util.emptyArray;
        /**
         * GetConversationsResponse total.
         * @member {number} total
         * @memberof messagespb.GetConversationsResponse
         * @instance
         */
        GetConversationsResponse.prototype.total = 0;
        /**
         * GetConversationsResponse page.
         * @member {number} page
         * @memberof messagespb.GetConversationsResponse
         * @instance
         */
        GetConversationsResponse.prototype.page = 0;
        /**
         * GetConversationsResponse limit.
         * @member {number} limit
         * @memberof messagespb.GetConversationsResponse
         * @instance
         */
        GetConversationsResponse.prototype.limit = 0;
        /**
         * Creates a new GetConversationsResponse instance using the specified properties.
         * @function create
         * @memberof messagespb.GetConversationsResponse
         * @static
         * @param {messagespb.IGetConversationsResponse=} [properties] Properties to set
         * @returns {messagespb.GetConversationsResponse} GetConversationsResponse instance
         */
        GetConversationsResponse.create = function create(properties) {
            return new GetConversationsResponse(properties);
        };
        /**
         * Encodes the specified GetConversationsResponse message. Does not implicitly {@link messagespb.GetConversationsResponse.verify|verify} messages.
         * @function encode
         * @memberof messagespb.GetConversationsResponse
         * @static
         * @param {messagespb.IGetConversationsResponse} message GetConversationsResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetConversationsResponse.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.conversations != null && message.conversations.length)
                for (let i = 0; i < message.conversations.length; ++i)
                    $root.messagespb.Conversation.encode(message.conversations[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.total != null && Object.hasOwnProperty.call(message, "total"))
                writer.uint32(/* id 2, wireType 0 =*/16).int32(message.total);
            if (message.page != null && Object.hasOwnProperty.call(message, "page"))
                writer.uint32(/* id 3, wireType 0 =*/24).int32(message.page);
            if (message.limit != null && Object.hasOwnProperty.call(message, "limit"))
                writer.uint32(/* id 4, wireType 0 =*/32).int32(message.limit);
            return writer;
        };
        /**
         * Encodes the specified GetConversationsResponse message, length delimited. Does not implicitly {@link messagespb.GetConversationsResponse.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.GetConversationsResponse
         * @static
         * @param {messagespb.IGetConversationsResponse} message GetConversationsResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetConversationsResponse.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a GetConversationsResponse message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.GetConversationsResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.GetConversationsResponse} GetConversationsResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetConversationsResponse.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.GetConversationsResponse();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        if (!(message.conversations && message.conversations.length))
                            message.conversations = [];
                        message.conversations.push($root.messagespb.Conversation.decode(reader, reader.uint32()));
                        break;
                    }
                case 2: {
                        message.total = reader.int32();
                        break;
                    }
                case 3: {
                        message.page = reader.int32();
                        break;
                    }
                case 4: {
                        message.limit = reader.int32();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a GetConversationsResponse message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.GetConversationsResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.GetConversationsResponse} GetConversationsResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetConversationsResponse.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a GetConversationsResponse message.
         * @function verify
         * @memberof messagespb.GetConversationsResponse
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        GetConversationsResponse.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.conversations != null && message.hasOwnProperty("conversations")) {
                if (!Array.isArray(message.conversations))
                    return "conversations: array expected";
                for (let i = 0; i < message.conversations.length; ++i) {
                    let error = $root.messagespb.Conversation.verify(message.conversations[i]);
                    if (error)
                        return "conversations." + error;
                }
            }
            if (message.total != null && message.hasOwnProperty("total"))
                if (!$util.isInteger(message.total))
                    return "total: integer expected";
            if (message.page != null && message.hasOwnProperty("page"))
                if (!$util.isInteger(message.page))
                    return "page: integer expected";
            if (message.limit != null && message.hasOwnProperty("limit"))
                if (!$util.isInteger(message.limit))
                    return "limit: integer expected";
            return null;
        };
        /**
         * Creates a GetConversationsResponse message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.GetConversationsResponse
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.GetConversationsResponse} GetConversationsResponse
         */
        GetConversationsResponse.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.GetConversationsResponse)
                return object;
            let message = new $root.messagespb.GetConversationsResponse();
            if (object.conversations) {
                if (!Array.isArray(object.conversations))
                    throw TypeError(".messagespb.GetConversationsResponse.conversations: array expected");
                message.conversations = [];
                for (let i = 0; i < object.conversations.length; ++i) {
                    if (typeof object.conversations[i] !== "object")
                        throw TypeError(".messagespb.GetConversationsResponse.conversations: object expected");
                    message.conversations[i] = $root.messagespb.Conversation.fromObject(object.conversations[i]);
                }
            }
            if (object.total != null)
                message.total = object.total | 0;
            if (object.page != null)
                message.page = object.page | 0;
            if (object.limit != null)
                message.limit = object.limit | 0;
            return message;
        };
        /**
         * Creates a plain object from a GetConversationsResponse message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.GetConversationsResponse
         * @static
         * @param {messagespb.GetConversationsResponse} message GetConversationsResponse
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        GetConversationsResponse.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.arrays || options.defaults)
                object.conversations = [];
            if (options.defaults) {
                object.total = 0;
                object.page = 0;
                object.limit = 0;
            }
            if (message.conversations && message.conversations.length) {
                object.conversations = [];
                for (let j = 0; j < message.conversations.length; ++j)
                    object.conversations[j] = $root.messagespb.Conversation.toObject(message.conversations[j], options);
            }
            if (message.total != null && message.hasOwnProperty("total"))
                object.total = message.total;
            if (message.page != null && message.hasOwnProperty("page"))
                object.page = message.page;
            if (message.limit != null && message.hasOwnProperty("limit"))
                object.limit = message.limit;
            return object;
        };
        /**
         * Converts this GetConversationsResponse to JSON.
         * @function toJSON
         * @memberof messagespb.GetConversationsResponse
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        GetConversationsResponse.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for GetConversationsResponse
         * @function getTypeUrl
         * @memberof messagespb.GetConversationsResponse
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        GetConversationsResponse.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.GetConversationsResponse";
        };
        return GetConversationsResponse;
    })();
    messagespb.GetActiveConversationsRequest = (function() {
        /**
         * Properties of a GetActiveConversationsRequest.
         * @memberof messagespb
         * @interface IGetActiveConversationsRequest
         * @property {string|null} [userId] GetActiveConversationsRequest userId
         * @property {number|null} [page] GetActiveConversationsRequest page
         * @property {number|null} [limit] GetActiveConversationsRequest limit
         */
        /**
         * Constructs a new GetActiveConversationsRequest.
         * @memberof messagespb
         * @classdesc Represents a GetActiveConversationsRequest.
         * @implements IGetActiveConversationsRequest
         * @constructor
         * @param {messagespb.IGetActiveConversationsRequest=} [properties] Properties to set
         */
        function GetActiveConversationsRequest(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * GetActiveConversationsRequest userId.
         * @member {string} userId
         * @memberof messagespb.GetActiveConversationsRequest
         * @instance
         */
        GetActiveConversationsRequest.prototype.userId = "";
        /**
         * GetActiveConversationsRequest page.
         * @member {number} page
         * @memberof messagespb.GetActiveConversationsRequest
         * @instance
         */
        GetActiveConversationsRequest.prototype.page = 0;
        /**
         * GetActiveConversationsRequest limit.
         * @member {number} limit
         * @memberof messagespb.GetActiveConversationsRequest
         * @instance
         */
        GetActiveConversationsRequest.prototype.limit = 0;
        /**
         * Creates a new GetActiveConversationsRequest instance using the specified properties.
         * @function create
         * @memberof messagespb.GetActiveConversationsRequest
         * @static
         * @param {messagespb.IGetActiveConversationsRequest=} [properties] Properties to set
         * @returns {messagespb.GetActiveConversationsRequest} GetActiveConversationsRequest instance
         */
        GetActiveConversationsRequest.create = function create(properties) {
            return new GetActiveConversationsRequest(properties);
        };
        /**
         * Encodes the specified GetActiveConversationsRequest message. Does not implicitly {@link messagespb.GetActiveConversationsRequest.verify|verify} messages.
         * @function encode
         * @memberof messagespb.GetActiveConversationsRequest
         * @static
         * @param {messagespb.IGetActiveConversationsRequest} message GetActiveConversationsRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetActiveConversationsRequest.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.userId != null && Object.hasOwnProperty.call(message, "userId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.userId);
            if (message.page != null && Object.hasOwnProperty.call(message, "page"))
                writer.uint32(/* id 2, wireType 0 =*/16).int32(message.page);
            if (message.limit != null && Object.hasOwnProperty.call(message, "limit"))
                writer.uint32(/* id 3, wireType 0 =*/24).int32(message.limit);
            return writer;
        };
        /**
         * Encodes the specified GetActiveConversationsRequest message, length delimited. Does not implicitly {@link messagespb.GetActiveConversationsRequest.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.GetActiveConversationsRequest
         * @static
         * @param {messagespb.IGetActiveConversationsRequest} message GetActiveConversationsRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetActiveConversationsRequest.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a GetActiveConversationsRequest message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.GetActiveConversationsRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.GetActiveConversationsRequest} GetActiveConversationsRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetActiveConversationsRequest.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.GetActiveConversationsRequest();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.userId = reader.string();
                        break;
                    }
                case 2: {
                        message.page = reader.int32();
                        break;
                    }
                case 3: {
                        message.limit = reader.int32();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a GetActiveConversationsRequest message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.GetActiveConversationsRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.GetActiveConversationsRequest} GetActiveConversationsRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetActiveConversationsRequest.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a GetActiveConversationsRequest message.
         * @function verify
         * @memberof messagespb.GetActiveConversationsRequest
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        GetActiveConversationsRequest.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.userId != null && message.hasOwnProperty("userId"))
                if (!$util.isString(message.userId))
                    return "userId: string expected";
            if (message.page != null && message.hasOwnProperty("page"))
                if (!$util.isInteger(message.page))
                    return "page: integer expected";
            if (message.limit != null && message.hasOwnProperty("limit"))
                if (!$util.isInteger(message.limit))
                    return "limit: integer expected";
            return null;
        };
        /**
         * Creates a GetActiveConversationsRequest message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.GetActiveConversationsRequest
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.GetActiveConversationsRequest} GetActiveConversationsRequest
         */
        GetActiveConversationsRequest.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.GetActiveConversationsRequest)
                return object;
            let message = new $root.messagespb.GetActiveConversationsRequest();
            if (object.userId != null)
                message.userId = String(object.userId);
            if (object.page != null)
                message.page = object.page | 0;
            if (object.limit != null)
                message.limit = object.limit | 0;
            return message;
        };
        /**
         * Creates a plain object from a GetActiveConversationsRequest message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.GetActiveConversationsRequest
         * @static
         * @param {messagespb.GetActiveConversationsRequest} message GetActiveConversationsRequest
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        GetActiveConversationsRequest.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.userId = "";
                object.page = 0;
                object.limit = 0;
            }
            if (message.userId != null && message.hasOwnProperty("userId"))
                object.userId = message.userId;
            if (message.page != null && message.hasOwnProperty("page"))
                object.page = message.page;
            if (message.limit != null && message.hasOwnProperty("limit"))
                object.limit = message.limit;
            return object;
        };
        /**
         * Converts this GetActiveConversationsRequest to JSON.
         * @function toJSON
         * @memberof messagespb.GetActiveConversationsRequest
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        GetActiveConversationsRequest.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for GetActiveConversationsRequest
         * @function getTypeUrl
         * @memberof messagespb.GetActiveConversationsRequest
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        GetActiveConversationsRequest.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.GetActiveConversationsRequest";
        };
        return GetActiveConversationsRequest;
    })();
    messagespb.GetActiveConversationsResponse = (function() {
        /**
         * Properties of a GetActiveConversationsResponse.
         * @memberof messagespb
         * @interface IGetActiveConversationsResponse
         * @property {Array.<messagespb.IConversation>|null} [conversations] GetActiveConversationsResponse conversations
         * @property {number|null} [total] GetActiveConversationsResponse total
         * @property {number|null} [page] GetActiveConversationsResponse page
         * @property {number|null} [limit] GetActiveConversationsResponse limit
         */
        /**
         * Constructs a new GetActiveConversationsResponse.
         * @memberof messagespb
         * @classdesc Represents a GetActiveConversationsResponse.
         * @implements IGetActiveConversationsResponse
         * @constructor
         * @param {messagespb.IGetActiveConversationsResponse=} [properties] Properties to set
         */
        function GetActiveConversationsResponse(properties) {
            this.conversations = [];
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * GetActiveConversationsResponse conversations.
         * @member {Array.<messagespb.IConversation>} conversations
         * @memberof messagespb.GetActiveConversationsResponse
         * @instance
         */
        GetActiveConversationsResponse.prototype.conversations = $util.emptyArray;
        /**
         * GetActiveConversationsResponse total.
         * @member {number} total
         * @memberof messagespb.GetActiveConversationsResponse
         * @instance
         */
        GetActiveConversationsResponse.prototype.total = 0;
        /**
         * GetActiveConversationsResponse page.
         * @member {number} page
         * @memberof messagespb.GetActiveConversationsResponse
         * @instance
         */
        GetActiveConversationsResponse.prototype.page = 0;
        /**
         * GetActiveConversationsResponse limit.
         * @member {number} limit
         * @memberof messagespb.GetActiveConversationsResponse
         * @instance
         */
        GetActiveConversationsResponse.prototype.limit = 0;
        /**
         * Creates a new GetActiveConversationsResponse instance using the specified properties.
         * @function create
         * @memberof messagespb.GetActiveConversationsResponse
         * @static
         * @param {messagespb.IGetActiveConversationsResponse=} [properties] Properties to set
         * @returns {messagespb.GetActiveConversationsResponse} GetActiveConversationsResponse instance
         */
        GetActiveConversationsResponse.create = function create(properties) {
            return new GetActiveConversationsResponse(properties);
        };
        /**
         * Encodes the specified GetActiveConversationsResponse message. Does not implicitly {@link messagespb.GetActiveConversationsResponse.verify|verify} messages.
         * @function encode
         * @memberof messagespb.GetActiveConversationsResponse
         * @static
         * @param {messagespb.IGetActiveConversationsResponse} message GetActiveConversationsResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetActiveConversationsResponse.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.conversations != null && message.conversations.length)
                for (let i = 0; i < message.conversations.length; ++i)
                    $root.messagespb.Conversation.encode(message.conversations[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.total != null && Object.hasOwnProperty.call(message, "total"))
                writer.uint32(/* id 2, wireType 0 =*/16).int32(message.total);
            if (message.page != null && Object.hasOwnProperty.call(message, "page"))
                writer.uint32(/* id 3, wireType 0 =*/24).int32(message.page);
            if (message.limit != null && Object.hasOwnProperty.call(message, "limit"))
                writer.uint32(/* id 4, wireType 0 =*/32).int32(message.limit);
            return writer;
        };
        /**
         * Encodes the specified GetActiveConversationsResponse message, length delimited. Does not implicitly {@link messagespb.GetActiveConversationsResponse.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.GetActiveConversationsResponse
         * @static
         * @param {messagespb.IGetActiveConversationsResponse} message GetActiveConversationsResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetActiveConversationsResponse.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a GetActiveConversationsResponse message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.GetActiveConversationsResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.GetActiveConversationsResponse} GetActiveConversationsResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetActiveConversationsResponse.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.GetActiveConversationsResponse();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        if (!(message.conversations && message.conversations.length))
                            message.conversations = [];
                        message.conversations.push($root.messagespb.Conversation.decode(reader, reader.uint32()));
                        break;
                    }
                case 2: {
                        message.total = reader.int32();
                        break;
                    }
                case 3: {
                        message.page = reader.int32();
                        break;
                    }
                case 4: {
                        message.limit = reader.int32();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a GetActiveConversationsResponse message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.GetActiveConversationsResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.GetActiveConversationsResponse} GetActiveConversationsResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetActiveConversationsResponse.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a GetActiveConversationsResponse message.
         * @function verify
         * @memberof messagespb.GetActiveConversationsResponse
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        GetActiveConversationsResponse.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.conversations != null && message.hasOwnProperty("conversations")) {
                if (!Array.isArray(message.conversations))
                    return "conversations: array expected";
                for (let i = 0; i < message.conversations.length; ++i) {
                    let error = $root.messagespb.Conversation.verify(message.conversations[i]);
                    if (error)
                        return "conversations." + error;
                }
            }
            if (message.total != null && message.hasOwnProperty("total"))
                if (!$util.isInteger(message.total))
                    return "total: integer expected";
            if (message.page != null && message.hasOwnProperty("page"))
                if (!$util.isInteger(message.page))
                    return "page: integer expected";
            if (message.limit != null && message.hasOwnProperty("limit"))
                if (!$util.isInteger(message.limit))
                    return "limit: integer expected";
            return null;
        };
        /**
         * Creates a GetActiveConversationsResponse message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.GetActiveConversationsResponse
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.GetActiveConversationsResponse} GetActiveConversationsResponse
         */
        GetActiveConversationsResponse.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.GetActiveConversationsResponse)
                return object;
            let message = new $root.messagespb.GetActiveConversationsResponse();
            if (object.conversations) {
                if (!Array.isArray(object.conversations))
                    throw TypeError(".messagespb.GetActiveConversationsResponse.conversations: array expected");
                message.conversations = [];
                for (let i = 0; i < object.conversations.length; ++i) {
                    if (typeof object.conversations[i] !== "object")
                        throw TypeError(".messagespb.GetActiveConversationsResponse.conversations: object expected");
                    message.conversations[i] = $root.messagespb.Conversation.fromObject(object.conversations[i]);
                }
            }
            if (object.total != null)
                message.total = object.total | 0;
            if (object.page != null)
                message.page = object.page | 0;
            if (object.limit != null)
                message.limit = object.limit | 0;
            return message;
        };
        /**
         * Creates a plain object from a GetActiveConversationsResponse message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.GetActiveConversationsResponse
         * @static
         * @param {messagespb.GetActiveConversationsResponse} message GetActiveConversationsResponse
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        GetActiveConversationsResponse.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.arrays || options.defaults)
                object.conversations = [];
            if (options.defaults) {
                object.total = 0;
                object.page = 0;
                object.limit = 0;
            }
            if (message.conversations && message.conversations.length) {
                object.conversations = [];
                for (let j = 0; j < message.conversations.length; ++j)
                    object.conversations[j] = $root.messagespb.Conversation.toObject(message.conversations[j], options);
            }
            if (message.total != null && message.hasOwnProperty("total"))
                object.total = message.total;
            if (message.page != null && message.hasOwnProperty("page"))
                object.page = message.page;
            if (message.limit != null && message.hasOwnProperty("limit"))
                object.limit = message.limit;
            return object;
        };
        /**
         * Converts this GetActiveConversationsResponse to JSON.
         * @function toJSON
         * @memberof messagespb.GetActiveConversationsResponse
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        GetActiveConversationsResponse.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for GetActiveConversationsResponse
         * @function getTypeUrl
         * @memberof messagespb.GetActiveConversationsResponse
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        GetActiveConversationsResponse.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.GetActiveConversationsResponse";
        };
        return GetActiveConversationsResponse;
    })();
    messagespb.SendMessageRequest = (function() {
        /**
         * Properties of a SendMessageRequest.
         * @memberof messagespb
         * @interface ISendMessageRequest
         * @property {string|null} [conversationId] SendMessageRequest conversationId
         * @property {string|null} [senderId] SendMessageRequest senderId
         * @property {string|null} [recipientId] SendMessageRequest recipientId
         * @property {string|null} [itemId] SendMessageRequest itemId
         * @property {string|null} [body] SendMessageRequest body
         */
        /**
         * Constructs a new SendMessageRequest.
         * @memberof messagespb
         * @classdesc Represents a SendMessageRequest.
         * @implements ISendMessageRequest
         * @constructor
         * @param {messagespb.ISendMessageRequest=} [properties] Properties to set
         */
        function SendMessageRequest(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * SendMessageRequest conversationId.
         * @member {string} conversationId
         * @memberof messagespb.SendMessageRequest
         * @instance
         */
        SendMessageRequest.prototype.conversationId = "";
        /**
         * SendMessageRequest senderId.
         * @member {string} senderId
         * @memberof messagespb.SendMessageRequest
         * @instance
         */
        SendMessageRequest.prototype.senderId = "";
        /**
         * SendMessageRequest recipientId.
         * @member {string} recipientId
         * @memberof messagespb.SendMessageRequest
         * @instance
         */
        SendMessageRequest.prototype.recipientId = "";
        /**
         * SendMessageRequest itemId.
         * @member {string} itemId
         * @memberof messagespb.SendMessageRequest
         * @instance
         */
        SendMessageRequest.prototype.itemId = "";
        /**
         * SendMessageRequest body.
         * @member {string} body
         * @memberof messagespb.SendMessageRequest
         * @instance
         */
        SendMessageRequest.prototype.body = "";
        /**
         * Creates a new SendMessageRequest instance using the specified properties.
         * @function create
         * @memberof messagespb.SendMessageRequest
         * @static
         * @param {messagespb.ISendMessageRequest=} [properties] Properties to set
         * @returns {messagespb.SendMessageRequest} SendMessageRequest instance
         */
        SendMessageRequest.create = function create(properties) {
            return new SendMessageRequest(properties);
        };
        /**
         * Encodes the specified SendMessageRequest message. Does not implicitly {@link messagespb.SendMessageRequest.verify|verify} messages.
         * @function encode
         * @memberof messagespb.SendMessageRequest
         * @static
         * @param {messagespb.ISendMessageRequest} message SendMessageRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        SendMessageRequest.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.conversationId != null && Object.hasOwnProperty.call(message, "conversationId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.conversationId);
            if (message.senderId != null && Object.hasOwnProperty.call(message, "senderId"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.senderId);
            if (message.recipientId != null && Object.hasOwnProperty.call(message, "recipientId"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.recipientId);
            if (message.itemId != null && Object.hasOwnProperty.call(message, "itemId"))
                writer.uint32(/* id 4, wireType 2 =*/34).string(message.itemId);
            if (message.body != null && Object.hasOwnProperty.call(message, "body"))
                writer.uint32(/* id 5, wireType 2 =*/42).string(message.body);
            return writer;
        };
        /**
         * Encodes the specified SendMessageRequest message, length delimited. Does not implicitly {@link messagespb.SendMessageRequest.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.SendMessageRequest
         * @static
         * @param {messagespb.ISendMessageRequest} message SendMessageRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        SendMessageRequest.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a SendMessageRequest message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.SendMessageRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.SendMessageRequest} SendMessageRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        SendMessageRequest.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.SendMessageRequest();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.conversationId = reader.string();
                        break;
                    }
                case 2: {
                        message.senderId = reader.string();
                        break;
                    }
                case 3: {
                        message.recipientId = reader.string();
                        break;
                    }
                case 4: {
                        message.itemId = reader.string();
                        break;
                    }
                case 5: {
                        message.body = reader.string();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a SendMessageRequest message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.SendMessageRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.SendMessageRequest} SendMessageRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        SendMessageRequest.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a SendMessageRequest message.
         * @function verify
         * @memberof messagespb.SendMessageRequest
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        SendMessageRequest.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.conversationId != null && message.hasOwnProperty("conversationId"))
                if (!$util.isString(message.conversationId))
                    return "conversationId: string expected";
            if (message.senderId != null && message.hasOwnProperty("senderId"))
                if (!$util.isString(message.senderId))
                    return "senderId: string expected";
            if (message.recipientId != null && message.hasOwnProperty("recipientId"))
                if (!$util.isString(message.recipientId))
                    return "recipientId: string expected";
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                if (!$util.isString(message.itemId))
                    return "itemId: string expected";
            if (message.body != null && message.hasOwnProperty("body"))
                if (!$util.isString(message.body))
                    return "body: string expected";
            return null;
        };
        /**
         * Creates a SendMessageRequest message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.SendMessageRequest
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.SendMessageRequest} SendMessageRequest
         */
        SendMessageRequest.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.SendMessageRequest)
                return object;
            let message = new $root.messagespb.SendMessageRequest();
            if (object.conversationId != null)
                message.conversationId = String(object.conversationId);
            if (object.senderId != null)
                message.senderId = String(object.senderId);
            if (object.recipientId != null)
                message.recipientId = String(object.recipientId);
            if (object.itemId != null)
                message.itemId = String(object.itemId);
            if (object.body != null)
                message.body = String(object.body);
            return message;
        };
        /**
         * Creates a plain object from a SendMessageRequest message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.SendMessageRequest
         * @static
         * @param {messagespb.SendMessageRequest} message SendMessageRequest
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        SendMessageRequest.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.conversationId = "";
                object.senderId = "";
                object.recipientId = "";
                object.itemId = "";
                object.body = "";
            }
            if (message.conversationId != null && message.hasOwnProperty("conversationId"))
                object.conversationId = message.conversationId;
            if (message.senderId != null && message.hasOwnProperty("senderId"))
                object.senderId = message.senderId;
            if (message.recipientId != null && message.hasOwnProperty("recipientId"))
                object.recipientId = message.recipientId;
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                object.itemId = message.itemId;
            if (message.body != null && message.hasOwnProperty("body"))
                object.body = message.body;
            return object;
        };
        /**
         * Converts this SendMessageRequest to JSON.
         * @function toJSON
         * @memberof messagespb.SendMessageRequest
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        SendMessageRequest.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for SendMessageRequest
         * @function getTypeUrl
         * @memberof messagespb.SendMessageRequest
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        SendMessageRequest.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.SendMessageRequest";
        };
        return SendMessageRequest;
    })();
    messagespb.SendMessageResponse = (function() {
        /**
         * Properties of a SendMessageResponse.
         * @memberof messagespb
         * @interface ISendMessageResponse
         * @property {string|null} [id] SendMessageResponse id
         * @property {google.protobuf.ITimestamp|null} [sentAt] SendMessageResponse sentAt
         */
        /**
         * Constructs a new SendMessageResponse.
         * @memberof messagespb
         * @classdesc Represents a SendMessageResponse.
         * @implements ISendMessageResponse
         * @constructor
         * @param {messagespb.ISendMessageResponse=} [properties] Properties to set
         */
        function SendMessageResponse(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * SendMessageResponse id.
         * @member {string} id
         * @memberof messagespb.SendMessageResponse
         * @instance
         */
        SendMessageResponse.prototype.id = "";
        /**
         * SendMessageResponse sentAt.
         * @member {google.protobuf.ITimestamp|null|undefined} sentAt
         * @memberof messagespb.SendMessageResponse
         * @instance
         */
        SendMessageResponse.prototype.sentAt = null;
        /**
         * Creates a new SendMessageResponse instance using the specified properties.
         * @function create
         * @memberof messagespb.SendMessageResponse
         * @static
         * @param {messagespb.ISendMessageResponse=} [properties] Properties to set
         * @returns {messagespb.SendMessageResponse} SendMessageResponse instance
         */
        SendMessageResponse.create = function create(properties) {
            return new SendMessageResponse(properties);
        };
        /**
         * Encodes the specified SendMessageResponse message. Does not implicitly {@link messagespb.SendMessageResponse.verify|verify} messages.
         * @function encode
         * @memberof messagespb.SendMessageResponse
         * @static
         * @param {messagespb.ISendMessageResponse} message SendMessageResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        SendMessageResponse.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
            if (message.sentAt != null && Object.hasOwnProperty.call(message, "sentAt"))
                $root.google.protobuf.Timestamp.encode(message.sentAt, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
            return writer;
        };
        /**
         * Encodes the specified SendMessageResponse message, length delimited. Does not implicitly {@link messagespb.SendMessageResponse.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.SendMessageResponse
         * @static
         * @param {messagespb.ISendMessageResponse} message SendMessageResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        SendMessageResponse.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a SendMessageResponse message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.SendMessageResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.SendMessageResponse} SendMessageResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        SendMessageResponse.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.SendMessageResponse();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.id = reader.string();
                        break;
                    }
                case 2: {
                        message.sentAt = $root.google.protobuf.Timestamp.decode(reader, reader.uint32());
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a SendMessageResponse message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.SendMessageResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.SendMessageResponse} SendMessageResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        SendMessageResponse.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a SendMessageResponse message.
         * @function verify
         * @memberof messagespb.SendMessageResponse
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        SendMessageResponse.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.id != null && message.hasOwnProperty("id"))
                if (!$util.isString(message.id))
                    return "id: string expected";
            if (message.sentAt != null && message.hasOwnProperty("sentAt")) {
                let error = $root.google.protobuf.Timestamp.verify(message.sentAt);
                if (error)
                    return "sentAt." + error;
            }
            return null;
        };
        /**
         * Creates a SendMessageResponse message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.SendMessageResponse
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.SendMessageResponse} SendMessageResponse
         */
        SendMessageResponse.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.SendMessageResponse)
                return object;
            let message = new $root.messagespb.SendMessageResponse();
            if (object.id != null)
                message.id = String(object.id);
            if (object.sentAt != null) {
                if (typeof object.sentAt !== "object")
                    throw TypeError(".messagespb.SendMessageResponse.sentAt: object expected");
                message.sentAt = $root.google.protobuf.Timestamp.fromObject(object.sentAt);
            }
            return message;
        };
        /**
         * Creates a plain object from a SendMessageResponse message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.SendMessageResponse
         * @static
         * @param {messagespb.SendMessageResponse} message SendMessageResponse
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        SendMessageResponse.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.id = "";
                object.sentAt = null;
            }
            if (message.id != null && message.hasOwnProperty("id"))
                object.id = message.id;
            if (message.sentAt != null && message.hasOwnProperty("sentAt"))
                object.sentAt = $root.google.protobuf.Timestamp.toObject(message.sentAt, options);
            return object;
        };
        /**
         * Converts this SendMessageResponse to JSON.
         * @function toJSON
         * @memberof messagespb.SendMessageResponse
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        SendMessageResponse.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for SendMessageResponse
         * @function getTypeUrl
         * @memberof messagespb.SendMessageResponse
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        SendMessageResponse.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.SendMessageResponse";
        };
        return SendMessageResponse;
    })();
    messagespb.DeleteMessageRequest = (function() {
        /**
         * Properties of a DeleteMessageRequest.
         * @memberof messagespb
         * @interface IDeleteMessageRequest
         * @property {string|null} [id] DeleteMessageRequest id
         */
        /**
         * Constructs a new DeleteMessageRequest.
         * @memberof messagespb
         * @classdesc Represents a DeleteMessageRequest.
         * @implements IDeleteMessageRequest
         * @constructor
         * @param {messagespb.IDeleteMessageRequest=} [properties] Properties to set
         */
        function DeleteMessageRequest(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * DeleteMessageRequest id.
         * @member {string} id
         * @memberof messagespb.DeleteMessageRequest
         * @instance
         */
        DeleteMessageRequest.prototype.id = "";
        /**
         * Creates a new DeleteMessageRequest instance using the specified properties.
         * @function create
         * @memberof messagespb.DeleteMessageRequest
         * @static
         * @param {messagespb.IDeleteMessageRequest=} [properties] Properties to set
         * @returns {messagespb.DeleteMessageRequest} DeleteMessageRequest instance
         */
        DeleteMessageRequest.create = function create(properties) {
            return new DeleteMessageRequest(properties);
        };
        /**
         * Encodes the specified DeleteMessageRequest message. Does not implicitly {@link messagespb.DeleteMessageRequest.verify|verify} messages.
         * @function encode
         * @memberof messagespb.DeleteMessageRequest
         * @static
         * @param {messagespb.IDeleteMessageRequest} message DeleteMessageRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        DeleteMessageRequest.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
            return writer;
        };
        /**
         * Encodes the specified DeleteMessageRequest message, length delimited. Does not implicitly {@link messagespb.DeleteMessageRequest.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.DeleteMessageRequest
         * @static
         * @param {messagespb.IDeleteMessageRequest} message DeleteMessageRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        DeleteMessageRequest.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a DeleteMessageRequest message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.DeleteMessageRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.DeleteMessageRequest} DeleteMessageRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        DeleteMessageRequest.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.DeleteMessageRequest();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.id = reader.string();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a DeleteMessageRequest message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.DeleteMessageRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.DeleteMessageRequest} DeleteMessageRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        DeleteMessageRequest.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a DeleteMessageRequest message.
         * @function verify
         * @memberof messagespb.DeleteMessageRequest
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        DeleteMessageRequest.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.id != null && message.hasOwnProperty("id"))
                if (!$util.isString(message.id))
                    return "id: string expected";
            return null;
        };
        /**
         * Creates a DeleteMessageRequest message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.DeleteMessageRequest
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.DeleteMessageRequest} DeleteMessageRequest
         */
        DeleteMessageRequest.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.DeleteMessageRequest)
                return object;
            let message = new $root.messagespb.DeleteMessageRequest();
            if (object.id != null)
                message.id = String(object.id);
            return message;
        };
        /**
         * Creates a plain object from a DeleteMessageRequest message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.DeleteMessageRequest
         * @static
         * @param {messagespb.DeleteMessageRequest} message DeleteMessageRequest
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        DeleteMessageRequest.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults)
                object.id = "";
            if (message.id != null && message.hasOwnProperty("id"))
                object.id = message.id;
            return object;
        };
        /**
         * Converts this DeleteMessageRequest to JSON.
         * @function toJSON
         * @memberof messagespb.DeleteMessageRequest
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        DeleteMessageRequest.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for DeleteMessageRequest
         * @function getTypeUrl
         * @memberof messagespb.DeleteMessageRequest
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        DeleteMessageRequest.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.DeleteMessageRequest";
        };
        return DeleteMessageRequest;
    })();
    messagespb.GetMessageRequest = (function() {
        /**
         * Properties of a GetMessageRequest.
         * @memberof messagespb
         * @interface IGetMessageRequest
         * @property {string|null} [id] GetMessageRequest id
         */
        /**
         * Constructs a new GetMessageRequest.
         * @memberof messagespb
         * @classdesc Represents a GetMessageRequest.
         * @implements IGetMessageRequest
         * @constructor
         * @param {messagespb.IGetMessageRequest=} [properties] Properties to set
         */
        function GetMessageRequest(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * GetMessageRequest id.
         * @member {string} id
         * @memberof messagespb.GetMessageRequest
         * @instance
         */
        GetMessageRequest.prototype.id = "";
        /**
         * Creates a new GetMessageRequest instance using the specified properties.
         * @function create
         * @memberof messagespb.GetMessageRequest
         * @static
         * @param {messagespb.IGetMessageRequest=} [properties] Properties to set
         * @returns {messagespb.GetMessageRequest} GetMessageRequest instance
         */
        GetMessageRequest.create = function create(properties) {
            return new GetMessageRequest(properties);
        };
        /**
         * Encodes the specified GetMessageRequest message. Does not implicitly {@link messagespb.GetMessageRequest.verify|verify} messages.
         * @function encode
         * @memberof messagespb.GetMessageRequest
         * @static
         * @param {messagespb.IGetMessageRequest} message GetMessageRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetMessageRequest.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.id != null && Object.hasOwnProperty.call(message, "id"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.id);
            return writer;
        };
        /**
         * Encodes the specified GetMessageRequest message, length delimited. Does not implicitly {@link messagespb.GetMessageRequest.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.GetMessageRequest
         * @static
         * @param {messagespb.IGetMessageRequest} message GetMessageRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetMessageRequest.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a GetMessageRequest message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.GetMessageRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.GetMessageRequest} GetMessageRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetMessageRequest.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.GetMessageRequest();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.id = reader.string();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a GetMessageRequest message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.GetMessageRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.GetMessageRequest} GetMessageRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetMessageRequest.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a GetMessageRequest message.
         * @function verify
         * @memberof messagespb.GetMessageRequest
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        GetMessageRequest.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.id != null && message.hasOwnProperty("id"))
                if (!$util.isString(message.id))
                    return "id: string expected";
            return null;
        };
        /**
         * Creates a GetMessageRequest message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.GetMessageRequest
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.GetMessageRequest} GetMessageRequest
         */
        GetMessageRequest.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.GetMessageRequest)
                return object;
            let message = new $root.messagespb.GetMessageRequest();
            if (object.id != null)
                message.id = String(object.id);
            return message;
        };
        /**
         * Creates a plain object from a GetMessageRequest message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.GetMessageRequest
         * @static
         * @param {messagespb.GetMessageRequest} message GetMessageRequest
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        GetMessageRequest.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults)
                object.id = "";
            if (message.id != null && message.hasOwnProperty("id"))
                object.id = message.id;
            return object;
        };
        /**
         * Converts this GetMessageRequest to JSON.
         * @function toJSON
         * @memberof messagespb.GetMessageRequest
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        GetMessageRequest.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for GetMessageRequest
         * @function getTypeUrl
         * @memberof messagespb.GetMessageRequest
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        GetMessageRequest.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.GetMessageRequest";
        };
        return GetMessageRequest;
    })();
    messagespb.GetMessageResponse = (function() {
        /**
         * Properties of a GetMessageResponse.
         * @memberof messagespb
         * @interface IGetMessageResponse
         * @property {messagespb.IMessage|null} [message] GetMessageResponse message
         */
        /**
         * Constructs a new GetMessageResponse.
         * @memberof messagespb
         * @classdesc Represents a GetMessageResponse.
         * @implements IGetMessageResponse
         * @constructor
         * @param {messagespb.IGetMessageResponse=} [properties] Properties to set
         */
        function GetMessageResponse(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * GetMessageResponse message.
         * @member {messagespb.IMessage|null|undefined} message
         * @memberof messagespb.GetMessageResponse
         * @instance
         */
        GetMessageResponse.prototype.message = null;
        /**
         * Creates a new GetMessageResponse instance using the specified properties.
         * @function create
         * @memberof messagespb.GetMessageResponse
         * @static
         * @param {messagespb.IGetMessageResponse=} [properties] Properties to set
         * @returns {messagespb.GetMessageResponse} GetMessageResponse instance
         */
        GetMessageResponse.create = function create(properties) {
            return new GetMessageResponse(properties);
        };
        /**
         * Encodes the specified GetMessageResponse message. Does not implicitly {@link messagespb.GetMessageResponse.verify|verify} messages.
         * @function encode
         * @memberof messagespb.GetMessageResponse
         * @static
         * @param {messagespb.IGetMessageResponse} message GetMessageResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetMessageResponse.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.message != null && Object.hasOwnProperty.call(message, "message"))
                $root.messagespb.Message.encode(message.message, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            return writer;
        };
        /**
         * Encodes the specified GetMessageResponse message, length delimited. Does not implicitly {@link messagespb.GetMessageResponse.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.GetMessageResponse
         * @static
         * @param {messagespb.IGetMessageResponse} message GetMessageResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetMessageResponse.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a GetMessageResponse message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.GetMessageResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.GetMessageResponse} GetMessageResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetMessageResponse.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.GetMessageResponse();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.message = $root.messagespb.Message.decode(reader, reader.uint32());
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a GetMessageResponse message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.GetMessageResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.GetMessageResponse} GetMessageResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetMessageResponse.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a GetMessageResponse message.
         * @function verify
         * @memberof messagespb.GetMessageResponse
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        GetMessageResponse.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.message != null && message.hasOwnProperty("message")) {
                let error = $root.messagespb.Message.verify(message.message);
                if (error)
                    return "message." + error;
            }
            return null;
        };
        /**
         * Creates a GetMessageResponse message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.GetMessageResponse
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.GetMessageResponse} GetMessageResponse
         */
        GetMessageResponse.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.GetMessageResponse)
                return object;
            let message = new $root.messagespb.GetMessageResponse();
            if (object.message != null) {
                if (typeof object.message !== "object")
                    throw TypeError(".messagespb.GetMessageResponse.message: object expected");
                message.message = $root.messagespb.Message.fromObject(object.message);
            }
            return message;
        };
        /**
         * Creates a plain object from a GetMessageResponse message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.GetMessageResponse
         * @static
         * @param {messagespb.GetMessageResponse} message GetMessageResponse
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        GetMessageResponse.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults)
                object.message = null;
            if (message.message != null && message.hasOwnProperty("message"))
                object.message = $root.messagespb.Message.toObject(message.message, options);
            return object;
        };
        /**
         * Converts this GetMessageResponse to JSON.
         * @function toJSON
         * @memberof messagespb.GetMessageResponse
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        GetMessageResponse.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for GetMessageResponse
         * @function getTypeUrl
         * @memberof messagespb.GetMessageResponse
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        GetMessageResponse.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.GetMessageResponse";
        };
        return GetMessageResponse;
    })();
    messagespb.GetMessagesRequest = (function() {
        /**
         * Properties of a GetMessagesRequest.
         * @memberof messagespb
         * @interface IGetMessagesRequest
         * @property {string|null} [conversationId] GetMessagesRequest conversationId
         * @property {number|null} [page] GetMessagesRequest page
         * @property {number|null} [limit] GetMessagesRequest limit
         */
        /**
         * Constructs a new GetMessagesRequest.
         * @memberof messagespb
         * @classdesc Represents a GetMessagesRequest.
         * @implements IGetMessagesRequest
         * @constructor
         * @param {messagespb.IGetMessagesRequest=} [properties] Properties to set
         */
        function GetMessagesRequest(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * GetMessagesRequest conversationId.
         * @member {string} conversationId
         * @memberof messagespb.GetMessagesRequest
         * @instance
         */
        GetMessagesRequest.prototype.conversationId = "";
        /**
         * GetMessagesRequest page.
         * @member {number} page
         * @memberof messagespb.GetMessagesRequest
         * @instance
         */
        GetMessagesRequest.prototype.page = 0;
        /**
         * GetMessagesRequest limit.
         * @member {number} limit
         * @memberof messagespb.GetMessagesRequest
         * @instance
         */
        GetMessagesRequest.prototype.limit = 0;
        /**
         * Creates a new GetMessagesRequest instance using the specified properties.
         * @function create
         * @memberof messagespb.GetMessagesRequest
         * @static
         * @param {messagespb.IGetMessagesRequest=} [properties] Properties to set
         * @returns {messagespb.GetMessagesRequest} GetMessagesRequest instance
         */
        GetMessagesRequest.create = function create(properties) {
            return new GetMessagesRequest(properties);
        };
        /**
         * Encodes the specified GetMessagesRequest message. Does not implicitly {@link messagespb.GetMessagesRequest.verify|verify} messages.
         * @function encode
         * @memberof messagespb.GetMessagesRequest
         * @static
         * @param {messagespb.IGetMessagesRequest} message GetMessagesRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetMessagesRequest.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.conversationId != null && Object.hasOwnProperty.call(message, "conversationId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.conversationId);
            if (message.page != null && Object.hasOwnProperty.call(message, "page"))
                writer.uint32(/* id 2, wireType 0 =*/16).int32(message.page);
            if (message.limit != null && Object.hasOwnProperty.call(message, "limit"))
                writer.uint32(/* id 3, wireType 0 =*/24).int32(message.limit);
            return writer;
        };
        /**
         * Encodes the specified GetMessagesRequest message, length delimited. Does not implicitly {@link messagespb.GetMessagesRequest.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.GetMessagesRequest
         * @static
         * @param {messagespb.IGetMessagesRequest} message GetMessagesRequest message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetMessagesRequest.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a GetMessagesRequest message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.GetMessagesRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.GetMessagesRequest} GetMessagesRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetMessagesRequest.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.GetMessagesRequest();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.conversationId = reader.string();
                        break;
                    }
                case 2: {
                        message.page = reader.int32();
                        break;
                    }
                case 3: {
                        message.limit = reader.int32();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a GetMessagesRequest message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.GetMessagesRequest
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.GetMessagesRequest} GetMessagesRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetMessagesRequest.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a GetMessagesRequest message.
         * @function verify
         * @memberof messagespb.GetMessagesRequest
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        GetMessagesRequest.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.conversationId != null && message.hasOwnProperty("conversationId"))
                if (!$util.isString(message.conversationId))
                    return "conversationId: string expected";
            if (message.page != null && message.hasOwnProperty("page"))
                if (!$util.isInteger(message.page))
                    return "page: integer expected";
            if (message.limit != null && message.hasOwnProperty("limit"))
                if (!$util.isInteger(message.limit))
                    return "limit: integer expected";
            return null;
        };
        /**
         * Creates a GetMessagesRequest message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.GetMessagesRequest
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.GetMessagesRequest} GetMessagesRequest
         */
        GetMessagesRequest.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.GetMessagesRequest)
                return object;
            let message = new $root.messagespb.GetMessagesRequest();
            if (object.conversationId != null)
                message.conversationId = String(object.conversationId);
            if (object.page != null)
                message.page = object.page | 0;
            if (object.limit != null)
                message.limit = object.limit | 0;
            return message;
        };
        /**
         * Creates a plain object from a GetMessagesRequest message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.GetMessagesRequest
         * @static
         * @param {messagespb.GetMessagesRequest} message GetMessagesRequest
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        GetMessagesRequest.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.conversationId = "";
                object.page = 0;
                object.limit = 0;
            }
            if (message.conversationId != null && message.hasOwnProperty("conversationId"))
                object.conversationId = message.conversationId;
            if (message.page != null && message.hasOwnProperty("page"))
                object.page = message.page;
            if (message.limit != null && message.hasOwnProperty("limit"))
                object.limit = message.limit;
            return object;
        };
        /**
         * Converts this GetMessagesRequest to JSON.
         * @function toJSON
         * @memberof messagespb.GetMessagesRequest
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        GetMessagesRequest.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for GetMessagesRequest
         * @function getTypeUrl
         * @memberof messagespb.GetMessagesRequest
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        GetMessagesRequest.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.GetMessagesRequest";
        };
        return GetMessagesRequest;
    })();
    messagespb.GetMessagesResponse = (function() {
        /**
         * Properties of a GetMessagesResponse.
         * @memberof messagespb
         * @interface IGetMessagesResponse
         * @property {Array.<messagespb.IMessage>|null} [messages] GetMessagesResponse messages
         * @property {number|null} [total] GetMessagesResponse total
         * @property {number|null} [page] GetMessagesResponse page
         * @property {number|null} [limit] GetMessagesResponse limit
         */
        /**
         * Constructs a new GetMessagesResponse.
         * @memberof messagespb
         * @classdesc Represents a GetMessagesResponse.
         * @implements IGetMessagesResponse
         * @constructor
         * @param {messagespb.IGetMessagesResponse=} [properties] Properties to set
         */
        function GetMessagesResponse(properties) {
            this.messages = [];
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * GetMessagesResponse messages.
         * @member {Array.<messagespb.IMessage>} messages
         * @memberof messagespb.GetMessagesResponse
         * @instance
         */
        GetMessagesResponse.prototype.messages = $util.emptyArray;
        /**
         * GetMessagesResponse total.
         * @member {number} total
         * @memberof messagespb.GetMessagesResponse
         * @instance
         */
        GetMessagesResponse.prototype.total = 0;
        /**
         * GetMessagesResponse page.
         * @member {number} page
         * @memberof messagespb.GetMessagesResponse
         * @instance
         */
        GetMessagesResponse.prototype.page = 0;
        /**
         * GetMessagesResponse limit.
         * @member {number} limit
         * @memberof messagespb.GetMessagesResponse
         * @instance
         */
        GetMessagesResponse.prototype.limit = 0;
        /**
         * Creates a new GetMessagesResponse instance using the specified properties.
         * @function create
         * @memberof messagespb.GetMessagesResponse
         * @static
         * @param {messagespb.IGetMessagesResponse=} [properties] Properties to set
         * @returns {messagespb.GetMessagesResponse} GetMessagesResponse instance
         */
        GetMessagesResponse.create = function create(properties) {
            return new GetMessagesResponse(properties);
        };
        /**
         * Encodes the specified GetMessagesResponse message. Does not implicitly {@link messagespb.GetMessagesResponse.verify|verify} messages.
         * @function encode
         * @memberof messagespb.GetMessagesResponse
         * @static
         * @param {messagespb.IGetMessagesResponse} message GetMessagesResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetMessagesResponse.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.messages != null && message.messages.length)
                for (let i = 0; i < message.messages.length; ++i)
                    $root.messagespb.Message.encode(message.messages[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.total != null && Object.hasOwnProperty.call(message, "total"))
                writer.uint32(/* id 2, wireType 0 =*/16).int32(message.total);
            if (message.page != null && Object.hasOwnProperty.call(message, "page"))
                writer.uint32(/* id 3, wireType 0 =*/24).int32(message.page);
            if (message.limit != null && Object.hasOwnProperty.call(message, "limit"))
                writer.uint32(/* id 4, wireType 0 =*/32).int32(message.limit);
            return writer;
        };
        /**
         * Encodes the specified GetMessagesResponse message, length delimited. Does not implicitly {@link messagespb.GetMessagesResponse.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.GetMessagesResponse
         * @static
         * @param {messagespb.IGetMessagesResponse} message GetMessagesResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GetMessagesResponse.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes a GetMessagesResponse message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.GetMessagesResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.GetMessagesResponse} GetMessagesResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetMessagesResponse.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.GetMessagesResponse();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        if (!(message.messages && message.messages.length))
                            message.messages = [];
                        message.messages.push($root.messagespb.Message.decode(reader, reader.uint32()));
                        break;
                    }
                case 2: {
                        message.total = reader.int32();
                        break;
                    }
                case 3: {
                        message.page = reader.int32();
                        break;
                    }
                case 4: {
                        message.limit = reader.int32();
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes a GetMessagesResponse message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.GetMessagesResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.GetMessagesResponse} GetMessagesResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GetMessagesResponse.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies a GetMessagesResponse message.
         * @function verify
         * @memberof messagespb.GetMessagesResponse
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        GetMessagesResponse.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.messages != null && message.hasOwnProperty("messages")) {
                if (!Array.isArray(message.messages))
                    return "messages: array expected";
                for (let i = 0; i < message.messages.length; ++i) {
                    let error = $root.messagespb.Message.verify(message.messages[i]);
                    if (error)
                        return "messages." + error;
                }
            }
            if (message.total != null && message.hasOwnProperty("total"))
                if (!$util.isInteger(message.total))
                    return "total: integer expected";
            if (message.page != null && message.hasOwnProperty("page"))
                if (!$util.isInteger(message.page))
                    return "page: integer expected";
            if (message.limit != null && message.hasOwnProperty("limit"))
                if (!$util.isInteger(message.limit))
                    return "limit: integer expected";
            return null;
        };
        /**
         * Creates a GetMessagesResponse message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.GetMessagesResponse
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.GetMessagesResponse} GetMessagesResponse
         */
        GetMessagesResponse.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.GetMessagesResponse)
                return object;
            let message = new $root.messagespb.GetMessagesResponse();
            if (object.messages) {
                if (!Array.isArray(object.messages))
                    throw TypeError(".messagespb.GetMessagesResponse.messages: array expected");
                message.messages = [];
                for (let i = 0; i < object.messages.length; ++i) {
                    if (typeof object.messages[i] !== "object")
                        throw TypeError(".messagespb.GetMessagesResponse.messages: object expected");
                    message.messages[i] = $root.messagespb.Message.fromObject(object.messages[i]);
                }
            }
            if (object.total != null)
                message.total = object.total | 0;
            if (object.page != null)
                message.page = object.page | 0;
            if (object.limit != null)
                message.limit = object.limit | 0;
            return message;
        };
        /**
         * Creates a plain object from a GetMessagesResponse message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.GetMessagesResponse
         * @static
         * @param {messagespb.GetMessagesResponse} message GetMessagesResponse
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        GetMessagesResponse.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.arrays || options.defaults)
                object.messages = [];
            if (options.defaults) {
                object.total = 0;
                object.page = 0;
                object.limit = 0;
            }
            if (message.messages && message.messages.length) {
                object.messages = [];
                for (let j = 0; j < message.messages.length; ++j)
                    object.messages[j] = $root.messagespb.Message.toObject(message.messages[j], options);
            }
            if (message.total != null && message.hasOwnProperty("total"))
                object.total = message.total;
            if (message.page != null && message.hasOwnProperty("page"))
                object.page = message.page;
            if (message.limit != null && message.hasOwnProperty("limit"))
                object.limit = message.limit;
            return object;
        };
        /**
         * Converts this GetMessagesResponse to JSON.
         * @function toJSON
         * @memberof messagespb.GetMessagesResponse
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        GetMessagesResponse.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for GetMessagesResponse
         * @function getTypeUrl
         * @memberof messagespb.GetMessagesResponse
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        GetMessagesResponse.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.GetMessagesResponse";
        };
        return GetMessagesResponse;
    })();
    messagespb.ErrorResponse = (function() {
        /**
         * Properties of an ErrorResponse.
         * @memberof messagespb
         * @interface IErrorResponse
         * @property {number|null} [code] ErrorResponse code
         * @property {string|null} [message] ErrorResponse message
         * @property {Array.<string>|null} [details] ErrorResponse details
         */
        /**
         * Constructs a new ErrorResponse.
         * @memberof messagespb
         * @classdesc Represents an ErrorResponse.
         * @implements IErrorResponse
         * @constructor
         * @param {messagespb.IErrorResponse=} [properties] Properties to set
         */
        function ErrorResponse(properties) {
            this.details = [];
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }
        /**
         * ErrorResponse code.
         * @member {number} code
         * @memberof messagespb.ErrorResponse
         * @instance
         */
        ErrorResponse.prototype.code = 0;
        /**
         * ErrorResponse message.
         * @member {string} message
         * @memberof messagespb.ErrorResponse
         * @instance
         */
        ErrorResponse.prototype.message = "";
        /**
         * ErrorResponse details.
         * @member {Array.<string>} details
         * @memberof messagespb.ErrorResponse
         * @instance
         */
        ErrorResponse.prototype.details = $util.emptyArray;
        /**
         * Creates a new ErrorResponse instance using the specified properties.
         * @function create
         * @memberof messagespb.ErrorResponse
         * @static
         * @param {messagespb.IErrorResponse=} [properties] Properties to set
         * @returns {messagespb.ErrorResponse} ErrorResponse instance
         */
        ErrorResponse.create = function create(properties) {
            return new ErrorResponse(properties);
        };
        /**
         * Encodes the specified ErrorResponse message. Does not implicitly {@link messagespb.ErrorResponse.verify|verify} messages.
         * @function encode
         * @memberof messagespb.ErrorResponse
         * @static
         * @param {messagespb.IErrorResponse} message ErrorResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ErrorResponse.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.code != null && Object.hasOwnProperty.call(message, "code"))
                writer.uint32(/* id 1, wireType 0 =*/8).int32(message.code);
            if (message.message != null && Object.hasOwnProperty.call(message, "message"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.message);
            if (message.details != null && message.details.length)
                for (let i = 0; i < message.details.length; ++i)
                    writer.uint32(/* id 3, wireType 2 =*/26).string(message.details[i]);
            return writer;
        };
        /**
         * Encodes the specified ErrorResponse message, length delimited. Does not implicitly {@link messagespb.ErrorResponse.verify|verify} messages.
         * @function encodeDelimited
         * @memberof messagespb.ErrorResponse
         * @static
         * @param {messagespb.IErrorResponse} message ErrorResponse message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ErrorResponse.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };
        /**
         * Decodes an ErrorResponse message from the specified reader or buffer.
         * @function decode
         * @memberof messagespb.ErrorResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {messagespb.ErrorResponse} ErrorResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ErrorResponse.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.messagespb.ErrorResponse();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1: {
                        message.code = reader.int32();
                        break;
                    }
                case 2: {
                        message.message = reader.string();
                        break;
                    }
                case 3: {
                        if (!(message.details && message.details.length))
                            message.details = [];
                        message.details.push(reader.string());
                        break;
                    }
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };
        /**
         * Decodes an ErrorResponse message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof messagespb.ErrorResponse
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {messagespb.ErrorResponse} ErrorResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ErrorResponse.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };
        /**
         * Verifies an ErrorResponse message.
         * @function verify
         * @memberof messagespb.ErrorResponse
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        ErrorResponse.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.code != null && message.hasOwnProperty("code"))
                if (!$util.isInteger(message.code))
                    return "code: integer expected";
            if (message.message != null && message.hasOwnProperty("message"))
                if (!$util.isString(message.message))
                    return "message: string expected";
            if (message.details != null && message.hasOwnProperty("details")) {
                if (!Array.isArray(message.details))
                    return "details: array expected";
                for (let i = 0; i < message.details.length; ++i)
                    if (!$util.isString(message.details[i]))
                        return "details: string[] expected";
            }
            return null;
        };
        /**
         * Creates an ErrorResponse message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof messagespb.ErrorResponse
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {messagespb.ErrorResponse} ErrorResponse
         */
        ErrorResponse.fromObject = function fromObject(object) {
            if (object instanceof $root.messagespb.ErrorResponse)
                return object;
            let message = new $root.messagespb.ErrorResponse();
            if (object.code != null)
                message.code = object.code | 0;
            if (object.message != null)
                message.message = String(object.message);
            if (object.details) {
                if (!Array.isArray(object.details))
                    throw TypeError(".messagespb.ErrorResponse.details: array expected");
                message.details = [];
                for (let i = 0; i < object.details.length; ++i)
                    message.details[i] = String(object.details[i]);
            }
            return message;
        };
        /**
         * Creates a plain object from an ErrorResponse message. Also converts values to other types if specified.
         * @function toObject
         * @memberof messagespb.ErrorResponse
         * @static
         * @param {messagespb.ErrorResponse} message ErrorResponse
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        ErrorResponse.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.arrays || options.defaults)
                object.details = [];
            if (options.defaults) {
                object.code = 0;
                object.message = "";
            }
            if (message.code != null && message.hasOwnProperty("code"))
                object.code = message.code;
            if (message.message != null && message.hasOwnProperty("message"))
                object.message = message.message;
            if (message.details && message.details.length) {
                object.details = [];
                for (let j = 0; j < message.details.length; ++j)
                    object.details[j] = message.details[j];
            }
            return object;
        };
        /**
         * Converts this ErrorResponse to JSON.
         * @function toJSON
         * @memberof messagespb.ErrorResponse
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        ErrorResponse.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };
        /**
         * Gets the default type url for ErrorResponse
         * @function getTypeUrl
         * @memberof messagespb.ErrorResponse
         * @static
         * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns {string} The default type url
         */
        ErrorResponse.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
            if (typeUrlPrefix === undefined) {
                typeUrlPrefix = "type.googleapis.com";
            }
            return typeUrlPrefix + "/messagespb.ErrorResponse";
        };
        return ErrorResponse;
    })();
    return messagespb;
})();
export const google = $root.google = (() => {
    /**
     * Namespace google.
     * @exports google
     * @namespace
     */
    const google = {};
    google.protobuf = (function() {
        /**
         * Namespace protobuf.
         * @memberof google
         * @namespace
         */
        const protobuf = {};
        protobuf.Empty = (function() {
            /**
             * Properties of an Empty.
             * @memberof google.protobuf
             * @interface IEmpty
             */
            /**
             * Constructs a new Empty.
             * @memberof google.protobuf
             * @classdesc Represents an Empty.
             * @implements IEmpty
             * @constructor
             * @param {google.protobuf.IEmpty=} [properties] Properties to set
             */
            function Empty(properties) {
                if (properties)
                    for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }
            /**
             * Creates a new Empty instance using the specified properties.
             * @function create
             * @memberof google.protobuf.Empty
             * @static
             * @param {google.protobuf.IEmpty=} [properties] Properties to set
             * @returns {google.protobuf.Empty} Empty instance
             */
            Empty.create = function create(properties) {
                return new Empty(properties);
            };
            /**
             * Encodes the specified Empty message. Does not implicitly {@link google.protobuf.Empty.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.Empty
             * @static
             * @param {google.protobuf.IEmpty} message Empty message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Empty.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                return writer;
            };
            /**
             * Encodes the specified Empty message, length delimited. Does not implicitly {@link google.protobuf.Empty.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.Empty
             * @static
             * @param {google.protobuf.IEmpty} message Empty message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Empty.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };
            /**
             * Decodes an Empty message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.Empty
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.Empty} Empty
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Empty.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                let end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.Empty();
                while (reader.pos < end) {
                    let tag = reader.uint32();
                    switch (tag >>> 3) {
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };
            /**
             * Decodes an Empty message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.Empty
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.Empty} Empty
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Empty.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };
            /**
             * Verifies an Empty message.
             * @function verify
             * @memberof google.protobuf.Empty
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            Empty.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                return null;
            };
            /**
             * Creates an Empty message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.Empty
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.Empty} Empty
             */
            Empty.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.Empty)
                    return object;
                return new $root.google.protobuf.Empty();
            };
            /**
             * Creates a plain object from an Empty message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.Empty
             * @static
             * @param {google.protobuf.Empty} message Empty
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            Empty.toObject = function toObject() {
                return {};
            };
            /**
             * Converts this Empty to JSON.
             * @function toJSON
             * @memberof google.protobuf.Empty
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            Empty.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };
            /**
             * Gets the default type url for Empty
             * @function getTypeUrl
             * @memberof google.protobuf.Empty
             * @static
             * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns {string} The default type url
             */
            Empty.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
                if (typeUrlPrefix === undefined) {
                    typeUrlPrefix = "type.googleapis.com";
                }
                return typeUrlPrefix + "/google.protobuf.Empty";
            };
            return Empty;
        })();
        protobuf.Timestamp = (function() {
            /**
             * Properties of a Timestamp.
             * @memberof google.protobuf
             * @interface ITimestamp
             * @property {number|Long|null} [seconds] Timestamp seconds
             * @property {number|null} [nanos] Timestamp nanos
             */
            /**
             * Constructs a new Timestamp.
             * @memberof google.protobuf
             * @classdesc Represents a Timestamp.
             * @implements ITimestamp
             * @constructor
             * @param {google.protobuf.ITimestamp=} [properties] Properties to set
             */
            function Timestamp(properties) {
                if (properties)
                    for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }
            /**
             * Timestamp seconds.
             * @member {number|Long} seconds
             * @memberof google.protobuf.Timestamp
             * @instance
             */
            Timestamp.prototype.seconds = $util.Long ? $util.Long.fromBits(0,0,false) : 0;
            /**
             * Timestamp nanos.
             * @member {number} nanos
             * @memberof google.protobuf.Timestamp
             * @instance
             */
            Timestamp.prototype.nanos = 0;
            /**
             * Creates a new Timestamp instance using the specified properties.
             * @function create
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {google.protobuf.ITimestamp=} [properties] Properties to set
             * @returns {google.protobuf.Timestamp} Timestamp instance
             */
            Timestamp.create = function create(properties) {
                return new Timestamp(properties);
            };
            /**
             * Encodes the specified Timestamp message. Does not implicitly {@link google.protobuf.Timestamp.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {google.protobuf.ITimestamp} message Timestamp message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Timestamp.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.seconds != null && Object.hasOwnProperty.call(message, "seconds"))
                    writer.uint32(/* id 1, wireType 0 =*/8).int64(message.seconds);
                if (message.nanos != null && Object.hasOwnProperty.call(message, "nanos"))
                    writer.uint32(/* id 2, wireType 0 =*/16).int32(message.nanos);
                return writer;
            };
            /**
             * Encodes the specified Timestamp message, length delimited. Does not implicitly {@link google.protobuf.Timestamp.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {google.protobuf.ITimestamp} message Timestamp message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            Timestamp.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };
            /**
             * Decodes a Timestamp message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.Timestamp} Timestamp
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Timestamp.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                let end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.Timestamp();
                while (reader.pos < end) {
                    let tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1: {
                            message.seconds = reader.int64();
                            break;
                        }
                    case 2: {
                            message.nanos = reader.int32();
                            break;
                        }
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };
            /**
             * Decodes a Timestamp message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.Timestamp} Timestamp
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            Timestamp.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };
            /**
             * Verifies a Timestamp message.
             * @function verify
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            Timestamp.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.seconds != null && message.hasOwnProperty("seconds"))
                    if (!$util.isInteger(message.seconds) && !(message.seconds && $util.isInteger(message.seconds.low) && $util.isInteger(message.seconds.high)))
                        return "seconds: integer|Long expected";
                if (message.nanos != null && message.hasOwnProperty("nanos"))
                    if (!$util.isInteger(message.nanos))
                        return "nanos: integer expected";
                return null;
            };
            /**
             * Creates a Timestamp message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.Timestamp} Timestamp
             */
            Timestamp.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.Timestamp)
                    return object;
                let message = new $root.google.protobuf.Timestamp();
                if (object.seconds != null)
                    if ($util.Long)
                        (message.seconds = $util.Long.fromValue(object.seconds)).unsigned = false;
                    else if (typeof object.seconds === "string")
                        message.seconds = parseInt(object.seconds, 10);
                    else if (typeof object.seconds === "number")
                        message.seconds = object.seconds;
                    else if (typeof object.seconds === "object")
                        message.seconds = new $util.LongBits(object.seconds.low >>> 0, object.seconds.high >>> 0).toNumber();
                if (object.nanos != null)
                    message.nanos = object.nanos | 0;
                return message;
            };
            /**
             * Creates a plain object from a Timestamp message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {google.protobuf.Timestamp} message Timestamp
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            Timestamp.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                let object = {};
                if (options.defaults) {
                    if ($util.Long) {
                        let long = new $util.Long(0, 0, false);
                        object.seconds = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                    } else
                        object.seconds = options.longs === String ? "0" : 0;
                    object.nanos = 0;
                }
                if (message.seconds != null && message.hasOwnProperty("seconds"))
                    if (typeof message.seconds === "number")
                        object.seconds = options.longs === String ? String(message.seconds) : message.seconds;
                    else
                        object.seconds = options.longs === String ? $util.Long.prototype.toString.call(message.seconds) : options.longs === Number ? new $util.LongBits(message.seconds.low >>> 0, message.seconds.high >>> 0).toNumber() : message.seconds;
                if (message.nanos != null && message.hasOwnProperty("nanos"))
                    object.nanos = message.nanos;
                return object;
            };
            /**
             * Converts this Timestamp to JSON.
             * @function toJSON
             * @memberof google.protobuf.Timestamp
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            Timestamp.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };
            /**
             * Gets the default type url for Timestamp
             * @function getTypeUrl
             * @memberof google.protobuf.Timestamp
             * @static
             * @param {string} [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns {string} The default type url
             */
            Timestamp.getTypeUrl = function getTypeUrl(typeUrlPrefix) {
                if (typeUrlPrefix === undefined) {
                    typeUrlPrefix = "type.googleapis.com";
                }
                return typeUrlPrefix + "/google.protobuf.Timestamp";
            };
            return Timestamp;
        })();
        return protobuf;
    })();
    return google;
})();
export { $root as default };
