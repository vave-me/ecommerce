import * as $protobuf from "protobufjs";
import Long = require("long");
/** Namespace messagespb. */
export namespace messagespb {

    /** Represents a MessagesService */
    class MessagesService extends $protobuf.rpc.Service {

        /**
         * Constructs a new MessagesService service.
         * @param rpcImpl RPC implementation
         * @param [requestDelimited=false] Whether requests are length-delimited
         * @param [responseDelimited=false] Whether responses are length-delimited
         */
        constructor(rpcImpl: $protobuf.RPCImpl, requestDelimited?: boolean, responseDelimited?: boolean);

        /**
         * Creates new MessagesService service using the specified rpc implementation.
         * @param rpcImpl RPC implementation
         * @param [requestDelimited=false] Whether requests are length-delimited
         * @param [responseDelimited=false] Whether responses are length-delimited
         * @returns RPC service. Useful where requests and/or responses are streamed.
         */
        public static create(rpcImpl: $protobuf.RPCImpl, requestDelimited?: boolean, responseDelimited?: boolean): MessagesService;

        /**
         * Calls StartConversation.
         * @param request StartConversationRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and StartConversationResponse
         */
        public startConversation(request: messagespb.IStartConversationRequest, callback: messagespb.MessagesService.StartConversationCallback): void;

        /**
         * Calls StartConversation.
         * @param request StartConversationRequest message or plain object
         * @returns Promise
         */
        public startConversation(request: messagespb.IStartConversationRequest): Promise<messagespb.StartConversationResponse>;

        /**
         * Calls RestoreConversation.
         * @param request RestoreConversationRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and RestoreConversationResponse
         */
        public restoreConversation(request: messagespb.IRestoreConversationRequest, callback: messagespb.MessagesService.RestoreConversationCallback): void;

        /**
         * Calls RestoreConversation.
         * @param request RestoreConversationRequest message or plain object
         * @returns Promise
         */
        public restoreConversation(request: messagespb.IRestoreConversationRequest): Promise<messagespb.RestoreConversationResponse>;

        /**
         * Calls ArchiveConversation.
         * @param request ArchiveConversationRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and ArchiveConversationResponse
         */
        public archiveConversation(request: messagespb.IArchiveConversationRequest, callback: messagespb.MessagesService.ArchiveConversationCallback): void;

        /**
         * Calls ArchiveConversation.
         * @param request ArchiveConversationRequest message or plain object
         * @returns Promise
         */
        public archiveConversation(request: messagespb.IArchiveConversationRequest): Promise<messagespb.ArchiveConversationResponse>;

        /**
         * Calls GetConversation.
         * @param request GetConversationRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and GetConversationResponse
         */
        public getConversation(request: messagespb.IGetConversationRequest, callback: messagespb.MessagesService.GetConversationCallback): void;

        /**
         * Calls GetConversation.
         * @param request GetConversationRequest message or plain object
         * @returns Promise
         */
        public getConversation(request: messagespb.IGetConversationRequest): Promise<messagespb.GetConversationResponse>;

        /**
         * Calls GetConversations.
         * @param request GetConversationsRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and GetConversationsResponse
         */
        public getConversations(request: messagespb.IGetConversationsRequest, callback: messagespb.MessagesService.GetConversationsCallback): void;

        /**
         * Calls GetConversations.
         * @param request GetConversationsRequest message or plain object
         * @returns Promise
         */
        public getConversations(request: messagespb.IGetConversationsRequest): Promise<messagespb.GetConversationsResponse>;

        /**
         * Calls GetActiveConversations.
         * @param request GetActiveConversationsRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and GetActiveConversationsResponse
         */
        public getActiveConversations(request: messagespb.IGetActiveConversationsRequest, callback: messagespb.MessagesService.GetActiveConversationsCallback): void;

        /**
         * Calls GetActiveConversations.
         * @param request GetActiveConversationsRequest message or plain object
         * @returns Promise
         */
        public getActiveConversations(request: messagespb.IGetActiveConversationsRequest): Promise<messagespb.GetActiveConversationsResponse>;

        /**
         * Calls SendMessage.
         * @param request SendMessageRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and SendMessageResponse
         */
        public sendMessage(request: messagespb.ISendMessageRequest, callback: messagespb.MessagesService.SendMessageCallback): void;

        /**
         * Calls SendMessage.
         * @param request SendMessageRequest message or plain object
         * @returns Promise
         */
        public sendMessage(request: messagespb.ISendMessageRequest): Promise<messagespb.SendMessageResponse>;

        /**
         * Calls DeleteMessage.
         * @param request DeleteMessageRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and DeleteMessageResponse
         */
        public deleteMessage(request: messagespb.IDeleteMessageRequest, callback: messagespb.MessagesService.DeleteMessageCallback): void;

        /**
         * Calls DeleteMessage.
         * @param request DeleteMessageRequest message or plain object
         * @returns Promise
         */
        public deleteMessage(request: messagespb.IDeleteMessageRequest): Promise<messagespb.DeleteMessageResponse>;

        /**
         * Calls GetMessage.
         * @param request GetMessageRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and GetMessageResponse
         */
        public getMessage(request: messagespb.IGetMessageRequest, callback: messagespb.MessagesService.GetMessageCallback): void;

        /**
         * Calls GetMessage.
         * @param request GetMessageRequest message or plain object
         * @returns Promise
         */
        public getMessage(request: messagespb.IGetMessageRequest): Promise<messagespb.GetMessageResponse>;

        /**
         * Calls GetMessages.
         * @param request GetMessagesRequest message or plain object
         * @param callback Node-style callback called with the error, if any, and GetMessagesResponse
         */
        public getMessages(request: messagespb.IGetMessagesRequest, callback: messagespb.MessagesService.GetMessagesCallback): void;

        /**
         * Calls GetMessages.
         * @param request GetMessagesRequest message or plain object
         * @returns Promise
         */
        public getMessages(request: messagespb.IGetMessagesRequest): Promise<messagespb.GetMessagesResponse>;
    }

    namespace MessagesService {

        /**
         * Callback as used by {@link messagespb.MessagesService#startConversation}.
         * @param error Error, if any
         * @param [response] StartConversationResponse
         */
        type StartConversationCallback = (error: (Error|null), response?: messagespb.StartConversationResponse) => void;

        /**
         * Callback as used by {@link messagespb.MessagesService#restoreConversation}.
         * @param error Error, if any
         * @param [response] RestoreConversationResponse
         */
        type RestoreConversationCallback = (error: (Error|null), response?: messagespb.RestoreConversationResponse) => void;

        /**
         * Callback as used by {@link messagespb.MessagesService#archiveConversation}.
         * @param error Error, if any
         * @param [response] ArchiveConversationResponse
         */
        type ArchiveConversationCallback = (error: (Error|null), response?: messagespb.ArchiveConversationResponse) => void;

        /**
         * Callback as used by {@link messagespb.MessagesService#getConversation}.
         * @param error Error, if any
         * @param [response] GetConversationResponse
         */
        type GetConversationCallback = (error: (Error|null), response?: messagespb.GetConversationResponse) => void;

        /**
         * Callback as used by {@link messagespb.MessagesService#getConversations}.
         * @param error Error, if any
         * @param [response] GetConversationsResponse
         */
        type GetConversationsCallback = (error: (Error|null), response?: messagespb.GetConversationsResponse) => void;

        /**
         * Callback as used by {@link messagespb.MessagesService#getActiveConversations}.
         * @param error Error, if any
         * @param [response] GetActiveConversationsResponse
         */
        type GetActiveConversationsCallback = (error: (Error|null), response?: messagespb.GetActiveConversationsResponse) => void;

        /**
         * Callback as used by {@link messagespb.MessagesService#sendMessage}.
         * @param error Error, if any
         * @param [response] SendMessageResponse
         */
        type SendMessageCallback = (error: (Error|null), response?: messagespb.SendMessageResponse) => void;

        /**
         * Callback as used by {@link messagespb.MessagesService#deleteMessage}.
         * @param error Error, if any
         * @param [response] DeleteMessageResponse
         */
        type DeleteMessageCallback = (error: (Error|null), response?: messagespb.DeleteMessageResponse) => void;

        /**
         * Callback as used by {@link messagespb.MessagesService#getMessage}.
         * @param error Error, if any
         * @param [response] GetMessageResponse
         */
        type GetMessageCallback = (error: (Error|null), response?: messagespb.GetMessageResponse) => void;

        /**
         * Callback as used by {@link messagespb.MessagesService#getMessages}.
         * @param error Error, if any
         * @param [response] GetMessagesResponse
         */
        type GetMessagesCallback = (error: (Error|null), response?: messagespb.GetMessagesResponse) => void;
    }

    /** Properties of a DeleteMessageResponse. */
    interface IDeleteMessageResponse {
    }

    /** Represents a DeleteMessageResponse. */
    class DeleteMessageResponse implements IDeleteMessageResponse {

        /**
         * Constructs a new DeleteMessageResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IDeleteMessageResponse);

        /**
         * Creates a new DeleteMessageResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns DeleteMessageResponse instance
         */
        public static create(properties?: messagespb.IDeleteMessageResponse): messagespb.DeleteMessageResponse;

        /**
         * Encodes the specified DeleteMessageResponse message. Does not implicitly {@link messagespb.DeleteMessageResponse.verify|verify} messages.
         * @param message DeleteMessageResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IDeleteMessageResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified DeleteMessageResponse message, length delimited. Does not implicitly {@link messagespb.DeleteMessageResponse.verify|verify} messages.
         * @param message DeleteMessageResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IDeleteMessageResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a DeleteMessageResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns DeleteMessageResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.DeleteMessageResponse;

        /**
         * Decodes a DeleteMessageResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns DeleteMessageResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.DeleteMessageResponse;

        /**
         * Verifies a DeleteMessageResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a DeleteMessageResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns DeleteMessageResponse
         */
        public static fromObject(object: { [k: string]: any }): messagespb.DeleteMessageResponse;

        /**
         * Creates a plain object from a DeleteMessageResponse message. Also converts values to other types if specified.
         * @param message DeleteMessageResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.DeleteMessageResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this DeleteMessageResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for DeleteMessageResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a Conversation. */
    interface IConversation {

        /** Conversation id */
        id?: (string|null);

        /** Conversation senderId */
        senderId?: (string|null);

        /** Conversation recipientId */
        recipientId?: (string|null);

        /** Conversation itemId */
        itemId?: (string|null);

        /** Conversation conversationStatus */
        conversationStatus?: (string|null);
    }

    /** Represents a Conversation. */
    class Conversation implements IConversation {

        /**
         * Constructs a new Conversation.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IConversation);

        /** Conversation id. */
        public id: string;

        /** Conversation senderId. */
        public senderId: string;

        /** Conversation recipientId. */
        public recipientId: string;

        /** Conversation itemId. */
        public itemId: string;

        /** Conversation conversationStatus. */
        public conversationStatus: string;

        /**
         * Creates a new Conversation instance using the specified properties.
         * @param [properties] Properties to set
         * @returns Conversation instance
         */
        public static create(properties?: messagespb.IConversation): messagespb.Conversation;

        /**
         * Encodes the specified Conversation message. Does not implicitly {@link messagespb.Conversation.verify|verify} messages.
         * @param message Conversation message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IConversation, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified Conversation message, length delimited. Does not implicitly {@link messagespb.Conversation.verify|verify} messages.
         * @param message Conversation message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IConversation, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a Conversation message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns Conversation
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.Conversation;

        /**
         * Decodes a Conversation message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns Conversation
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.Conversation;

        /**
         * Verifies a Conversation message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a Conversation message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns Conversation
         */
        public static fromObject(object: { [k: string]: any }): messagespb.Conversation;

        /**
         * Creates a plain object from a Conversation message. Also converts values to other types if specified.
         * @param message Conversation
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.Conversation, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this Conversation to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for Conversation
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a Message. */
    interface IMessage {

        /** Message id */
        id?: (string|null);

        /** Message conversationId */
        conversationId?: (string|null);

        /** Message senderId */
        senderId?: (string|null);

        /** Message recipientId */
        recipientId?: (string|null);

        /** Message itemId */
        itemId?: (string|null);

        /** Message body */
        body?: (string|null);
    }

    /** Represents a Message. */
    class Message implements IMessage {

        /**
         * Constructs a new Message.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IMessage);

        /** Message id. */
        public id: string;

        /** Message conversationId. */
        public conversationId: string;

        /** Message senderId. */
        public senderId: string;

        /** Message recipientId. */
        public recipientId: string;

        /** Message itemId. */
        public itemId: string;

        /** Message body. */
        public body: string;

        /**
         * Creates a new Message instance using the specified properties.
         * @param [properties] Properties to set
         * @returns Message instance
         */
        public static create(properties?: messagespb.IMessage): messagespb.Message;

        /**
         * Encodes the specified Message message. Does not implicitly {@link messagespb.Message.verify|verify} messages.
         * @param message Message message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IMessage, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified Message message, length delimited. Does not implicitly {@link messagespb.Message.verify|verify} messages.
         * @param message Message message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IMessage, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a Message message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns Message
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.Message;

        /**
         * Decodes a Message message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns Message
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.Message;

        /**
         * Verifies a Message message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a Message message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns Message
         */
        public static fromObject(object: { [k: string]: any }): messagespb.Message;

        /**
         * Creates a plain object from a Message message. Also converts values to other types if specified.
         * @param message Message
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.Message, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this Message to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for Message
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** ConversationStatus enum. */
    enum ConversationStatus {
        PENDING = 0,
        ACTIVE = 1,
        ARCHIVED = 2
    }

    /** MessageStatus enum. */
    enum MessageStatus {
        SENT = 0,
        EDITED = 1,
        DELETED = 2
    }

    /** Properties of a StartConversationRequest. */
    interface IStartConversationRequest {

        /** StartConversationRequest senderId */
        senderId?: (string|null);

        /** StartConversationRequest recipientId */
        recipientId?: (string|null);

        /** StartConversationRequest itemId */
        itemId?: (string|null);
    }

    /** Represents a StartConversationRequest. */
    class StartConversationRequest implements IStartConversationRequest {

        /**
         * Constructs a new StartConversationRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IStartConversationRequest);

        /** StartConversationRequest senderId. */
        public senderId: string;

        /** StartConversationRequest recipientId. */
        public recipientId: string;

        /** StartConversationRequest itemId. */
        public itemId: string;

        /**
         * Creates a new StartConversationRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns StartConversationRequest instance
         */
        public static create(properties?: messagespb.IStartConversationRequest): messagespb.StartConversationRequest;

        /**
         * Encodes the specified StartConversationRequest message. Does not implicitly {@link messagespb.StartConversationRequest.verify|verify} messages.
         * @param message StartConversationRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IStartConversationRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified StartConversationRequest message, length delimited. Does not implicitly {@link messagespb.StartConversationRequest.verify|verify} messages.
         * @param message StartConversationRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IStartConversationRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a StartConversationRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns StartConversationRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.StartConversationRequest;

        /**
         * Decodes a StartConversationRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns StartConversationRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.StartConversationRequest;

        /**
         * Verifies a StartConversationRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a StartConversationRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns StartConversationRequest
         */
        public static fromObject(object: { [k: string]: any }): messagespb.StartConversationRequest;

        /**
         * Creates a plain object from a StartConversationRequest message. Also converts values to other types if specified.
         * @param message StartConversationRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.StartConversationRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this StartConversationRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for StartConversationRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a StartConversationResponse. */
    interface IStartConversationResponse {

        /** StartConversationResponse id */
        id?: (string|null);
    }

    /** Represents a StartConversationResponse. */
    class StartConversationResponse implements IStartConversationResponse {

        /**
         * Constructs a new StartConversationResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IStartConversationResponse);

        /** StartConversationResponse id. */
        public id: string;

        /**
         * Creates a new StartConversationResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns StartConversationResponse instance
         */
        public static create(properties?: messagespb.IStartConversationResponse): messagespb.StartConversationResponse;

        /**
         * Encodes the specified StartConversationResponse message. Does not implicitly {@link messagespb.StartConversationResponse.verify|verify} messages.
         * @param message StartConversationResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IStartConversationResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified StartConversationResponse message, length delimited. Does not implicitly {@link messagespb.StartConversationResponse.verify|verify} messages.
         * @param message StartConversationResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IStartConversationResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a StartConversationResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns StartConversationResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.StartConversationResponse;

        /**
         * Decodes a StartConversationResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns StartConversationResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.StartConversationResponse;

        /**
         * Verifies a StartConversationResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a StartConversationResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns StartConversationResponse
         */
        public static fromObject(object: { [k: string]: any }): messagespb.StartConversationResponse;

        /**
         * Creates a plain object from a StartConversationResponse message. Also converts values to other types if specified.
         * @param message StartConversationResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.StartConversationResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this StartConversationResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for StartConversationResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a RestoreConversationRequest. */
    interface IRestoreConversationRequest {

        /** RestoreConversationRequest id */
        id?: (string|null);
    }

    /** Represents a RestoreConversationRequest. */
    class RestoreConversationRequest implements IRestoreConversationRequest {

        /**
         * Constructs a new RestoreConversationRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IRestoreConversationRequest);

        /** RestoreConversationRequest id. */
        public id: string;

        /**
         * Creates a new RestoreConversationRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns RestoreConversationRequest instance
         */
        public static create(properties?: messagespb.IRestoreConversationRequest): messagespb.RestoreConversationRequest;

        /**
         * Encodes the specified RestoreConversationRequest message. Does not implicitly {@link messagespb.RestoreConversationRequest.verify|verify} messages.
         * @param message RestoreConversationRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IRestoreConversationRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified RestoreConversationRequest message, length delimited. Does not implicitly {@link messagespb.RestoreConversationRequest.verify|verify} messages.
         * @param message RestoreConversationRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IRestoreConversationRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a RestoreConversationRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns RestoreConversationRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.RestoreConversationRequest;

        /**
         * Decodes a RestoreConversationRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns RestoreConversationRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.RestoreConversationRequest;

        /**
         * Verifies a RestoreConversationRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a RestoreConversationRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns RestoreConversationRequest
         */
        public static fromObject(object: { [k: string]: any }): messagespb.RestoreConversationRequest;

        /**
         * Creates a plain object from a RestoreConversationRequest message. Also converts values to other types if specified.
         * @param message RestoreConversationRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.RestoreConversationRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this RestoreConversationRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for RestoreConversationRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a RestoreConversationResponse. */
    interface IRestoreConversationResponse {

        /** RestoreConversationResponse id */
        id?: (string|null);

        /** RestoreConversationResponse conversationStatus */
        conversationStatus?: (string|null);
    }

    /** Represents a RestoreConversationResponse. */
    class RestoreConversationResponse implements IRestoreConversationResponse {

        /**
         * Constructs a new RestoreConversationResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IRestoreConversationResponse);

        /** RestoreConversationResponse id. */
        public id: string;

        /** RestoreConversationResponse conversationStatus. */
        public conversationStatus: string;

        /**
         * Creates a new RestoreConversationResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns RestoreConversationResponse instance
         */
        public static create(properties?: messagespb.IRestoreConversationResponse): messagespb.RestoreConversationResponse;

        /**
         * Encodes the specified RestoreConversationResponse message. Does not implicitly {@link messagespb.RestoreConversationResponse.verify|verify} messages.
         * @param message RestoreConversationResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IRestoreConversationResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified RestoreConversationResponse message, length delimited. Does not implicitly {@link messagespb.RestoreConversationResponse.verify|verify} messages.
         * @param message RestoreConversationResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IRestoreConversationResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a RestoreConversationResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns RestoreConversationResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.RestoreConversationResponse;

        /**
         * Decodes a RestoreConversationResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns RestoreConversationResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.RestoreConversationResponse;

        /**
         * Verifies a RestoreConversationResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a RestoreConversationResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns RestoreConversationResponse
         */
        public static fromObject(object: { [k: string]: any }): messagespb.RestoreConversationResponse;

        /**
         * Creates a plain object from a RestoreConversationResponse message. Also converts values to other types if specified.
         * @param message RestoreConversationResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.RestoreConversationResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this RestoreConversationResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for RestoreConversationResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an ArchiveConversationRequest. */
    interface IArchiveConversationRequest {

        /** ArchiveConversationRequest id */
        id?: (string|null);
    }

    /** Represents an ArchiveConversationRequest. */
    class ArchiveConversationRequest implements IArchiveConversationRequest {

        /**
         * Constructs a new ArchiveConversationRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IArchiveConversationRequest);

        /** ArchiveConversationRequest id. */
        public id: string;

        /**
         * Creates a new ArchiveConversationRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns ArchiveConversationRequest instance
         */
        public static create(properties?: messagespb.IArchiveConversationRequest): messagespb.ArchiveConversationRequest;

        /**
         * Encodes the specified ArchiveConversationRequest message. Does not implicitly {@link messagespb.ArchiveConversationRequest.verify|verify} messages.
         * @param message ArchiveConversationRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IArchiveConversationRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified ArchiveConversationRequest message, length delimited. Does not implicitly {@link messagespb.ArchiveConversationRequest.verify|verify} messages.
         * @param message ArchiveConversationRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IArchiveConversationRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes an ArchiveConversationRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns ArchiveConversationRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.ArchiveConversationRequest;

        /**
         * Decodes an ArchiveConversationRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns ArchiveConversationRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.ArchiveConversationRequest;

        /**
         * Verifies an ArchiveConversationRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates an ArchiveConversationRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns ArchiveConversationRequest
         */
        public static fromObject(object: { [k: string]: any }): messagespb.ArchiveConversationRequest;

        /**
         * Creates a plain object from an ArchiveConversationRequest message. Also converts values to other types if specified.
         * @param message ArchiveConversationRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.ArchiveConversationRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this ArchiveConversationRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for ArchiveConversationRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an ArchiveConversationResponse. */
    interface IArchiveConversationResponse {

        /** ArchiveConversationResponse id */
        id?: (string|null);

        /** ArchiveConversationResponse conversationStatus */
        conversationStatus?: (string|null);
    }

    /** Represents an ArchiveConversationResponse. */
    class ArchiveConversationResponse implements IArchiveConversationResponse {

        /**
         * Constructs a new ArchiveConversationResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IArchiveConversationResponse);

        /** ArchiveConversationResponse id. */
        public id: string;

        /** ArchiveConversationResponse conversationStatus. */
        public conversationStatus: string;

        /**
         * Creates a new ArchiveConversationResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns ArchiveConversationResponse instance
         */
        public static create(properties?: messagespb.IArchiveConversationResponse): messagespb.ArchiveConversationResponse;

        /**
         * Encodes the specified ArchiveConversationResponse message. Does not implicitly {@link messagespb.ArchiveConversationResponse.verify|verify} messages.
         * @param message ArchiveConversationResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IArchiveConversationResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified ArchiveConversationResponse message, length delimited. Does not implicitly {@link messagespb.ArchiveConversationResponse.verify|verify} messages.
         * @param message ArchiveConversationResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IArchiveConversationResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes an ArchiveConversationResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns ArchiveConversationResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.ArchiveConversationResponse;

        /**
         * Decodes an ArchiveConversationResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns ArchiveConversationResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.ArchiveConversationResponse;

        /**
         * Verifies an ArchiveConversationResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates an ArchiveConversationResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns ArchiveConversationResponse
         */
        public static fromObject(object: { [k: string]: any }): messagespb.ArchiveConversationResponse;

        /**
         * Creates a plain object from an ArchiveConversationResponse message. Also converts values to other types if specified.
         * @param message ArchiveConversationResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.ArchiveConversationResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this ArchiveConversationResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for ArchiveConversationResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetConversationRequest. */
    interface IGetConversationRequest {

        /** GetConversationRequest id */
        id?: (string|null);
    }

    /** Represents a GetConversationRequest. */
    class GetConversationRequest implements IGetConversationRequest {

        /**
         * Constructs a new GetConversationRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IGetConversationRequest);

        /** GetConversationRequest id. */
        public id: string;

        /**
         * Creates a new GetConversationRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns GetConversationRequest instance
         */
        public static create(properties?: messagespb.IGetConversationRequest): messagespb.GetConversationRequest;

        /**
         * Encodes the specified GetConversationRequest message. Does not implicitly {@link messagespb.GetConversationRequest.verify|verify} messages.
         * @param message GetConversationRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IGetConversationRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified GetConversationRequest message, length delimited. Does not implicitly {@link messagespb.GetConversationRequest.verify|verify} messages.
         * @param message GetConversationRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IGetConversationRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a GetConversationRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns GetConversationRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.GetConversationRequest;

        /**
         * Decodes a GetConversationRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns GetConversationRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.GetConversationRequest;

        /**
         * Verifies a GetConversationRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a GetConversationRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns GetConversationRequest
         */
        public static fromObject(object: { [k: string]: any }): messagespb.GetConversationRequest;

        /**
         * Creates a plain object from a GetConversationRequest message. Also converts values to other types if specified.
         * @param message GetConversationRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.GetConversationRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this GetConversationRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for GetConversationRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetConversationResponse. */
    interface IGetConversationResponse {

        /** GetConversationResponse conversation */
        conversation?: (messagespb.IConversation|null);
    }

    /** Represents a GetConversationResponse. */
    class GetConversationResponse implements IGetConversationResponse {

        /**
         * Constructs a new GetConversationResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IGetConversationResponse);

        /** GetConversationResponse conversation. */
        public conversation?: (messagespb.IConversation|null);

        /**
         * Creates a new GetConversationResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns GetConversationResponse instance
         */
        public static create(properties?: messagespb.IGetConversationResponse): messagespb.GetConversationResponse;

        /**
         * Encodes the specified GetConversationResponse message. Does not implicitly {@link messagespb.GetConversationResponse.verify|verify} messages.
         * @param message GetConversationResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IGetConversationResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified GetConversationResponse message, length delimited. Does not implicitly {@link messagespb.GetConversationResponse.verify|verify} messages.
         * @param message GetConversationResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IGetConversationResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a GetConversationResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns GetConversationResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.GetConversationResponse;

        /**
         * Decodes a GetConversationResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns GetConversationResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.GetConversationResponse;

        /**
         * Verifies a GetConversationResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a GetConversationResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns GetConversationResponse
         */
        public static fromObject(object: { [k: string]: any }): messagespb.GetConversationResponse;

        /**
         * Creates a plain object from a GetConversationResponse message. Also converts values to other types if specified.
         * @param message GetConversationResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.GetConversationResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this GetConversationResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for GetConversationResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetConversationsRequest. */
    interface IGetConversationsRequest {

        /** GetConversationsRequest userId */
        userId?: (string|null);

        /** GetConversationsRequest page */
        page?: (number|null);

        /** GetConversationsRequest limit */
        limit?: (number|null);
    }

    /** Represents a GetConversationsRequest. */
    class GetConversationsRequest implements IGetConversationsRequest {

        /**
         * Constructs a new GetConversationsRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IGetConversationsRequest);

        /** GetConversationsRequest userId. */
        public userId: string;

        /** GetConversationsRequest page. */
        public page: number;

        /** GetConversationsRequest limit. */
        public limit: number;

        /**
         * Creates a new GetConversationsRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns GetConversationsRequest instance
         */
        public static create(properties?: messagespb.IGetConversationsRequest): messagespb.GetConversationsRequest;

        /**
         * Encodes the specified GetConversationsRequest message. Does not implicitly {@link messagespb.GetConversationsRequest.verify|verify} messages.
         * @param message GetConversationsRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IGetConversationsRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified GetConversationsRequest message, length delimited. Does not implicitly {@link messagespb.GetConversationsRequest.verify|verify} messages.
         * @param message GetConversationsRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IGetConversationsRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a GetConversationsRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns GetConversationsRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.GetConversationsRequest;

        /**
         * Decodes a GetConversationsRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns GetConversationsRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.GetConversationsRequest;

        /**
         * Verifies a GetConversationsRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a GetConversationsRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns GetConversationsRequest
         */
        public static fromObject(object: { [k: string]: any }): messagespb.GetConversationsRequest;

        /**
         * Creates a plain object from a GetConversationsRequest message. Also converts values to other types if specified.
         * @param message GetConversationsRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.GetConversationsRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this GetConversationsRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for GetConversationsRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetConversationsResponse. */
    interface IGetConversationsResponse {

        /** GetConversationsResponse conversations */
        conversations?: (messagespb.IConversation[]|null);

        /** GetConversationsResponse total */
        total?: (number|null);

        /** GetConversationsResponse page */
        page?: (number|null);

        /** GetConversationsResponse limit */
        limit?: (number|null);
    }

    /** Represents a GetConversationsResponse. */
    class GetConversationsResponse implements IGetConversationsResponse {

        /**
         * Constructs a new GetConversationsResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IGetConversationsResponse);

        /** GetConversationsResponse conversations. */
        public conversations: messagespb.IConversation[];

        /** GetConversationsResponse total. */
        public total: number;

        /** GetConversationsResponse page. */
        public page: number;

        /** GetConversationsResponse limit. */
        public limit: number;

        /**
         * Creates a new GetConversationsResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns GetConversationsResponse instance
         */
        public static create(properties?: messagespb.IGetConversationsResponse): messagespb.GetConversationsResponse;

        /**
         * Encodes the specified GetConversationsResponse message. Does not implicitly {@link messagespb.GetConversationsResponse.verify|verify} messages.
         * @param message GetConversationsResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IGetConversationsResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified GetConversationsResponse message, length delimited. Does not implicitly {@link messagespb.GetConversationsResponse.verify|verify} messages.
         * @param message GetConversationsResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IGetConversationsResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a GetConversationsResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns GetConversationsResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.GetConversationsResponse;

        /**
         * Decodes a GetConversationsResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns GetConversationsResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.GetConversationsResponse;

        /**
         * Verifies a GetConversationsResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a GetConversationsResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns GetConversationsResponse
         */
        public static fromObject(object: { [k: string]: any }): messagespb.GetConversationsResponse;

        /**
         * Creates a plain object from a GetConversationsResponse message. Also converts values to other types if specified.
         * @param message GetConversationsResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.GetConversationsResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this GetConversationsResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for GetConversationsResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetActiveConversationsRequest. */
    interface IGetActiveConversationsRequest {

        /** GetActiveConversationsRequest userId */
        userId?: (string|null);

        /** GetActiveConversationsRequest page */
        page?: (number|null);

        /** GetActiveConversationsRequest limit */
        limit?: (number|null);
    }

    /** Represents a GetActiveConversationsRequest. */
    class GetActiveConversationsRequest implements IGetActiveConversationsRequest {

        /**
         * Constructs a new GetActiveConversationsRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IGetActiveConversationsRequest);

        /** GetActiveConversationsRequest userId. */
        public userId: string;

        /** GetActiveConversationsRequest page. */
        public page: number;

        /** GetActiveConversationsRequest limit. */
        public limit: number;

        /**
         * Creates a new GetActiveConversationsRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns GetActiveConversationsRequest instance
         */
        public static create(properties?: messagespb.IGetActiveConversationsRequest): messagespb.GetActiveConversationsRequest;

        /**
         * Encodes the specified GetActiveConversationsRequest message. Does not implicitly {@link messagespb.GetActiveConversationsRequest.verify|verify} messages.
         * @param message GetActiveConversationsRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IGetActiveConversationsRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified GetActiveConversationsRequest message, length delimited. Does not implicitly {@link messagespb.GetActiveConversationsRequest.verify|verify} messages.
         * @param message GetActiveConversationsRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IGetActiveConversationsRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a GetActiveConversationsRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns GetActiveConversationsRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.GetActiveConversationsRequest;

        /**
         * Decodes a GetActiveConversationsRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns GetActiveConversationsRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.GetActiveConversationsRequest;

        /**
         * Verifies a GetActiveConversationsRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a GetActiveConversationsRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns GetActiveConversationsRequest
         */
        public static fromObject(object: { [k: string]: any }): messagespb.GetActiveConversationsRequest;

        /**
         * Creates a plain object from a GetActiveConversationsRequest message. Also converts values to other types if specified.
         * @param message GetActiveConversationsRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.GetActiveConversationsRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this GetActiveConversationsRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for GetActiveConversationsRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetActiveConversationsResponse. */
    interface IGetActiveConversationsResponse {

        /** GetActiveConversationsResponse conversations */
        conversations?: (messagespb.IConversation[]|null);

        /** GetActiveConversationsResponse total */
        total?: (number|null);

        /** GetActiveConversationsResponse page */
        page?: (number|null);

        /** GetActiveConversationsResponse limit */
        limit?: (number|null);
    }

    /** Represents a GetActiveConversationsResponse. */
    class GetActiveConversationsResponse implements IGetActiveConversationsResponse {

        /**
         * Constructs a new GetActiveConversationsResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IGetActiveConversationsResponse);

        /** GetActiveConversationsResponse conversations. */
        public conversations: messagespb.IConversation[];

        /** GetActiveConversationsResponse total. */
        public total: number;

        /** GetActiveConversationsResponse page. */
        public page: number;

        /** GetActiveConversationsResponse limit. */
        public limit: number;

        /**
         * Creates a new GetActiveConversationsResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns GetActiveConversationsResponse instance
         */
        public static create(properties?: messagespb.IGetActiveConversationsResponse): messagespb.GetActiveConversationsResponse;

        /**
         * Encodes the specified GetActiveConversationsResponse message. Does not implicitly {@link messagespb.GetActiveConversationsResponse.verify|verify} messages.
         * @param message GetActiveConversationsResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IGetActiveConversationsResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified GetActiveConversationsResponse message, length delimited. Does not implicitly {@link messagespb.GetActiveConversationsResponse.verify|verify} messages.
         * @param message GetActiveConversationsResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IGetActiveConversationsResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a GetActiveConversationsResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns GetActiveConversationsResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.GetActiveConversationsResponse;

        /**
         * Decodes a GetActiveConversationsResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns GetActiveConversationsResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.GetActiveConversationsResponse;

        /**
         * Verifies a GetActiveConversationsResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a GetActiveConversationsResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns GetActiveConversationsResponse
         */
        public static fromObject(object: { [k: string]: any }): messagespb.GetActiveConversationsResponse;

        /**
         * Creates a plain object from a GetActiveConversationsResponse message. Also converts values to other types if specified.
         * @param message GetActiveConversationsResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.GetActiveConversationsResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this GetActiveConversationsResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for GetActiveConversationsResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a SendMessageRequest. */
    interface ISendMessageRequest {

        /** SendMessageRequest conversationId */
        conversationId?: (string|null);

        /** SendMessageRequest senderId */
        senderId?: (string|null);

        /** SendMessageRequest recipientId */
        recipientId?: (string|null);

        /** SendMessageRequest itemId */
        itemId?: (string|null);

        /** SendMessageRequest body */
        body?: (string|null);
    }

    /** Represents a SendMessageRequest. */
    class SendMessageRequest implements ISendMessageRequest {

        /**
         * Constructs a new SendMessageRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.ISendMessageRequest);

        /** SendMessageRequest conversationId. */
        public conversationId: string;

        /** SendMessageRequest senderId. */
        public senderId: string;

        /** SendMessageRequest recipientId. */
        public recipientId: string;

        /** SendMessageRequest itemId. */
        public itemId: string;

        /** SendMessageRequest body. */
        public body: string;

        /**
         * Creates a new SendMessageRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns SendMessageRequest instance
         */
        public static create(properties?: messagespb.ISendMessageRequest): messagespb.SendMessageRequest;

        /**
         * Encodes the specified SendMessageRequest message. Does not implicitly {@link messagespb.SendMessageRequest.verify|verify} messages.
         * @param message SendMessageRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.ISendMessageRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified SendMessageRequest message, length delimited. Does not implicitly {@link messagespb.SendMessageRequest.verify|verify} messages.
         * @param message SendMessageRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.ISendMessageRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a SendMessageRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns SendMessageRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.SendMessageRequest;

        /**
         * Decodes a SendMessageRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns SendMessageRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.SendMessageRequest;

        /**
         * Verifies a SendMessageRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a SendMessageRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns SendMessageRequest
         */
        public static fromObject(object: { [k: string]: any }): messagespb.SendMessageRequest;

        /**
         * Creates a plain object from a SendMessageRequest message. Also converts values to other types if specified.
         * @param message SendMessageRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.SendMessageRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this SendMessageRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for SendMessageRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a SendMessageResponse. */
    interface ISendMessageResponse {

        /** SendMessageResponse id */
        id?: (string|null);

        /** SendMessageResponse sentAt */
        sentAt?: (google.protobuf.ITimestamp|null);
    }

    /** Represents a SendMessageResponse. */
    class SendMessageResponse implements ISendMessageResponse {

        /**
         * Constructs a new SendMessageResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.ISendMessageResponse);

        /** SendMessageResponse id. */
        public id: string;

        /** SendMessageResponse sentAt. */
        public sentAt?: (google.protobuf.ITimestamp|null);

        /**
         * Creates a new SendMessageResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns SendMessageResponse instance
         */
        public static create(properties?: messagespb.ISendMessageResponse): messagespb.SendMessageResponse;

        /**
         * Encodes the specified SendMessageResponse message. Does not implicitly {@link messagespb.SendMessageResponse.verify|verify} messages.
         * @param message SendMessageResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.ISendMessageResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified SendMessageResponse message, length delimited. Does not implicitly {@link messagespb.SendMessageResponse.verify|verify} messages.
         * @param message SendMessageResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.ISendMessageResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a SendMessageResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns SendMessageResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.SendMessageResponse;

        /**
         * Decodes a SendMessageResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns SendMessageResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.SendMessageResponse;

        /**
         * Verifies a SendMessageResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a SendMessageResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns SendMessageResponse
         */
        public static fromObject(object: { [k: string]: any }): messagespb.SendMessageResponse;

        /**
         * Creates a plain object from a SendMessageResponse message. Also converts values to other types if specified.
         * @param message SendMessageResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.SendMessageResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this SendMessageResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for SendMessageResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a DeleteMessageRequest. */
    interface IDeleteMessageRequest {

        /** DeleteMessageRequest id */
        id?: (string|null);
    }

    /** Represents a DeleteMessageRequest. */
    class DeleteMessageRequest implements IDeleteMessageRequest {

        /**
         * Constructs a new DeleteMessageRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IDeleteMessageRequest);

        /** DeleteMessageRequest id. */
        public id: string;

        /**
         * Creates a new DeleteMessageRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns DeleteMessageRequest instance
         */
        public static create(properties?: messagespb.IDeleteMessageRequest): messagespb.DeleteMessageRequest;

        /**
         * Encodes the specified DeleteMessageRequest message. Does not implicitly {@link messagespb.DeleteMessageRequest.verify|verify} messages.
         * @param message DeleteMessageRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IDeleteMessageRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified DeleteMessageRequest message, length delimited. Does not implicitly {@link messagespb.DeleteMessageRequest.verify|verify} messages.
         * @param message DeleteMessageRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IDeleteMessageRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a DeleteMessageRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns DeleteMessageRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.DeleteMessageRequest;

        /**
         * Decodes a DeleteMessageRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns DeleteMessageRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.DeleteMessageRequest;

        /**
         * Verifies a DeleteMessageRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a DeleteMessageRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns DeleteMessageRequest
         */
        public static fromObject(object: { [k: string]: any }): messagespb.DeleteMessageRequest;

        /**
         * Creates a plain object from a DeleteMessageRequest message. Also converts values to other types if specified.
         * @param message DeleteMessageRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.DeleteMessageRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this DeleteMessageRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for DeleteMessageRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetMessageRequest. */
    interface IGetMessageRequest {

        /** GetMessageRequest id */
        id?: (string|null);
    }

    /** Represents a GetMessageRequest. */
    class GetMessageRequest implements IGetMessageRequest {

        /**
         * Constructs a new GetMessageRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IGetMessageRequest);

        /** GetMessageRequest id. */
        public id: string;

        /**
         * Creates a new GetMessageRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns GetMessageRequest instance
         */
        public static create(properties?: messagespb.IGetMessageRequest): messagespb.GetMessageRequest;

        /**
         * Encodes the specified GetMessageRequest message. Does not implicitly {@link messagespb.GetMessageRequest.verify|verify} messages.
         * @param message GetMessageRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IGetMessageRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified GetMessageRequest message, length delimited. Does not implicitly {@link messagespb.GetMessageRequest.verify|verify} messages.
         * @param message GetMessageRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IGetMessageRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a GetMessageRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns GetMessageRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.GetMessageRequest;

        /**
         * Decodes a GetMessageRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns GetMessageRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.GetMessageRequest;

        /**
         * Verifies a GetMessageRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a GetMessageRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns GetMessageRequest
         */
        public static fromObject(object: { [k: string]: any }): messagespb.GetMessageRequest;

        /**
         * Creates a plain object from a GetMessageRequest message. Also converts values to other types if specified.
         * @param message GetMessageRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.GetMessageRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this GetMessageRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for GetMessageRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetMessageResponse. */
    interface IGetMessageResponse {

        /** GetMessageResponse message */
        message?: (messagespb.IMessage|null);
    }

    /** Represents a GetMessageResponse. */
    class GetMessageResponse implements IGetMessageResponse {

        /**
         * Constructs a new GetMessageResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IGetMessageResponse);

        /** GetMessageResponse message. */
        public message?: (messagespb.IMessage|null);

        /**
         * Creates a new GetMessageResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns GetMessageResponse instance
         */
        public static create(properties?: messagespb.IGetMessageResponse): messagespb.GetMessageResponse;

        /**
         * Encodes the specified GetMessageResponse message. Does not implicitly {@link messagespb.GetMessageResponse.verify|verify} messages.
         * @param message GetMessageResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IGetMessageResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified GetMessageResponse message, length delimited. Does not implicitly {@link messagespb.GetMessageResponse.verify|verify} messages.
         * @param message GetMessageResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IGetMessageResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a GetMessageResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns GetMessageResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.GetMessageResponse;

        /**
         * Decodes a GetMessageResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns GetMessageResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.GetMessageResponse;

        /**
         * Verifies a GetMessageResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a GetMessageResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns GetMessageResponse
         */
        public static fromObject(object: { [k: string]: any }): messagespb.GetMessageResponse;

        /**
         * Creates a plain object from a GetMessageResponse message. Also converts values to other types if specified.
         * @param message GetMessageResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.GetMessageResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this GetMessageResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for GetMessageResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetMessagesRequest. */
    interface IGetMessagesRequest {

        /** GetMessagesRequest conversationId */
        conversationId?: (string|null);

        /** GetMessagesRequest page */
        page?: (number|null);

        /** GetMessagesRequest limit */
        limit?: (number|null);
    }

    /** Represents a GetMessagesRequest. */
    class GetMessagesRequest implements IGetMessagesRequest {

        /**
         * Constructs a new GetMessagesRequest.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IGetMessagesRequest);

        /** GetMessagesRequest conversationId. */
        public conversationId: string;

        /** GetMessagesRequest page. */
        public page: number;

        /** GetMessagesRequest limit. */
        public limit: number;

        /**
         * Creates a new GetMessagesRequest instance using the specified properties.
         * @param [properties] Properties to set
         * @returns GetMessagesRequest instance
         */
        public static create(properties?: messagespb.IGetMessagesRequest): messagespb.GetMessagesRequest;

        /**
         * Encodes the specified GetMessagesRequest message. Does not implicitly {@link messagespb.GetMessagesRequest.verify|verify} messages.
         * @param message GetMessagesRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IGetMessagesRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified GetMessagesRequest message, length delimited. Does not implicitly {@link messagespb.GetMessagesRequest.verify|verify} messages.
         * @param message GetMessagesRequest message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IGetMessagesRequest, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a GetMessagesRequest message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns GetMessagesRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.GetMessagesRequest;

        /**
         * Decodes a GetMessagesRequest message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns GetMessagesRequest
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.GetMessagesRequest;

        /**
         * Verifies a GetMessagesRequest message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a GetMessagesRequest message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns GetMessagesRequest
         */
        public static fromObject(object: { [k: string]: any }): messagespb.GetMessagesRequest;

        /**
         * Creates a plain object from a GetMessagesRequest message. Also converts values to other types if specified.
         * @param message GetMessagesRequest
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.GetMessagesRequest, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this GetMessagesRequest to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for GetMessagesRequest
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of a GetMessagesResponse. */
    interface IGetMessagesResponse {

        /** GetMessagesResponse messages */
        messages?: (messagespb.IMessage[]|null);

        /** GetMessagesResponse total */
        total?: (number|null);

        /** GetMessagesResponse page */
        page?: (number|null);

        /** GetMessagesResponse limit */
        limit?: (number|null);
    }

    /** Represents a GetMessagesResponse. */
    class GetMessagesResponse implements IGetMessagesResponse {

        /**
         * Constructs a new GetMessagesResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IGetMessagesResponse);

        /** GetMessagesResponse messages. */
        public messages: messagespb.IMessage[];

        /** GetMessagesResponse total. */
        public total: number;

        /** GetMessagesResponse page. */
        public page: number;

        /** GetMessagesResponse limit. */
        public limit: number;

        /**
         * Creates a new GetMessagesResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns GetMessagesResponse instance
         */
        public static create(properties?: messagespb.IGetMessagesResponse): messagespb.GetMessagesResponse;

        /**
         * Encodes the specified GetMessagesResponse message. Does not implicitly {@link messagespb.GetMessagesResponse.verify|verify} messages.
         * @param message GetMessagesResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IGetMessagesResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified GetMessagesResponse message, length delimited. Does not implicitly {@link messagespb.GetMessagesResponse.verify|verify} messages.
         * @param message GetMessagesResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IGetMessagesResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes a GetMessagesResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns GetMessagesResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.GetMessagesResponse;

        /**
         * Decodes a GetMessagesResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns GetMessagesResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.GetMessagesResponse;

        /**
         * Verifies a GetMessagesResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates a GetMessagesResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns GetMessagesResponse
         */
        public static fromObject(object: { [k: string]: any }): messagespb.GetMessagesResponse;

        /**
         * Creates a plain object from a GetMessagesResponse message. Also converts values to other types if specified.
         * @param message GetMessagesResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.GetMessagesResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this GetMessagesResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for GetMessagesResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }

    /** Properties of an ErrorResponse. */
    interface IErrorResponse {

        /** ErrorResponse code */
        code?: (number|null);

        /** ErrorResponse message */
        message?: (string|null);

        /** ErrorResponse details */
        details?: (string[]|null);
    }

    /** Represents an ErrorResponse. */
    class ErrorResponse implements IErrorResponse {

        /**
         * Constructs a new ErrorResponse.
         * @param [properties] Properties to set
         */
        constructor(properties?: messagespb.IErrorResponse);

        /** ErrorResponse code. */
        public code: number;

        /** ErrorResponse message. */
        public message: string;

        /** ErrorResponse details. */
        public details: string[];

        /**
         * Creates a new ErrorResponse instance using the specified properties.
         * @param [properties] Properties to set
         * @returns ErrorResponse instance
         */
        public static create(properties?: messagespb.IErrorResponse): messagespb.ErrorResponse;

        /**
         * Encodes the specified ErrorResponse message. Does not implicitly {@link messagespb.ErrorResponse.verify|verify} messages.
         * @param message ErrorResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encode(message: messagespb.IErrorResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Encodes the specified ErrorResponse message, length delimited. Does not implicitly {@link messagespb.ErrorResponse.verify|verify} messages.
         * @param message ErrorResponse message or plain object to encode
         * @param [writer] Writer to encode to
         * @returns Writer
         */
        public static encodeDelimited(message: messagespb.IErrorResponse, writer?: $protobuf.Writer): $protobuf.Writer;

        /**
         * Decodes an ErrorResponse message from the specified reader or buffer.
         * @param reader Reader or buffer to decode from
         * @param [length] Message length if known beforehand
         * @returns ErrorResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): messagespb.ErrorResponse;

        /**
         * Decodes an ErrorResponse message from the specified reader or buffer, length delimited.
         * @param reader Reader or buffer to decode from
         * @returns ErrorResponse
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): messagespb.ErrorResponse;

        /**
         * Verifies an ErrorResponse message.
         * @param message Plain object to verify
         * @returns `null` if valid, otherwise the reason why it is not
         */
        public static verify(message: { [k: string]: any }): (string|null);

        /**
         * Creates an ErrorResponse message from a plain object. Also converts values to their respective internal types.
         * @param object Plain object
         * @returns ErrorResponse
         */
        public static fromObject(object: { [k: string]: any }): messagespb.ErrorResponse;

        /**
         * Creates a plain object from an ErrorResponse message. Also converts values to other types if specified.
         * @param message ErrorResponse
         * @param [options] Conversion options
         * @returns Plain object
         */
        public static toObject(message: messagespb.ErrorResponse, options?: $protobuf.IConversionOptions): { [k: string]: any };

        /**
         * Converts this ErrorResponse to JSON.
         * @returns JSON object
         */
        public toJSON(): { [k: string]: any };

        /**
         * Gets the default type url for ErrorResponse
         * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
         * @returns The default type url
         */
        public static getTypeUrl(typeUrlPrefix?: string): string;
    }
}

/** Namespace google. */
export namespace google {

    /** Namespace protobuf. */
    namespace protobuf {

        /** Properties of an Empty. */
        interface IEmpty {
        }

        /** Represents an Empty. */
        class Empty implements IEmpty {

            /**
             * Constructs a new Empty.
             * @param [properties] Properties to set
             */
            constructor(properties?: google.protobuf.IEmpty);

            /**
             * Creates a new Empty instance using the specified properties.
             * @param [properties] Properties to set
             * @returns Empty instance
             */
            public static create(properties?: google.protobuf.IEmpty): google.protobuf.Empty;

            /**
             * Encodes the specified Empty message. Does not implicitly {@link google.protobuf.Empty.verify|verify} messages.
             * @param message Empty message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encode(message: google.protobuf.IEmpty, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Encodes the specified Empty message, length delimited. Does not implicitly {@link google.protobuf.Empty.verify|verify} messages.
             * @param message Empty message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encodeDelimited(message: google.protobuf.IEmpty, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Decodes an Empty message from the specified reader or buffer.
             * @param reader Reader or buffer to decode from
             * @param [length] Message length if known beforehand
             * @returns Empty
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): google.protobuf.Empty;

            /**
             * Decodes an Empty message from the specified reader or buffer, length delimited.
             * @param reader Reader or buffer to decode from
             * @returns Empty
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): google.protobuf.Empty;

            /**
             * Verifies an Empty message.
             * @param message Plain object to verify
             * @returns `null` if valid, otherwise the reason why it is not
             */
            public static verify(message: { [k: string]: any }): (string|null);

            /**
             * Creates an Empty message from a plain object. Also converts values to their respective internal types.
             * @param object Plain object
             * @returns Empty
             */
            public static fromObject(object: { [k: string]: any }): google.protobuf.Empty;

            /**
             * Creates a plain object from an Empty message. Also converts values to other types if specified.
             * @param message Empty
             * @param [options] Conversion options
             * @returns Plain object
             */
            public static toObject(message: google.protobuf.Empty, options?: $protobuf.IConversionOptions): { [k: string]: any };

            /**
             * Converts this Empty to JSON.
             * @returns JSON object
             */
            public toJSON(): { [k: string]: any };

            /**
             * Gets the default type url for Empty
             * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns The default type url
             */
            public static getTypeUrl(typeUrlPrefix?: string): string;
        }

        /** Properties of a Timestamp. */
        interface ITimestamp {

            /** Timestamp seconds */
            seconds?: (number|Long|null);

            /** Timestamp nanos */
            nanos?: (number|null);
        }

        /** Represents a Timestamp. */
        class Timestamp implements ITimestamp {

            /**
             * Constructs a new Timestamp.
             * @param [properties] Properties to set
             */
            constructor(properties?: google.protobuf.ITimestamp);

            /** Timestamp seconds. */
            public seconds: (number|Long);

            /** Timestamp nanos. */
            public nanos: number;

            /**
             * Creates a new Timestamp instance using the specified properties.
             * @param [properties] Properties to set
             * @returns Timestamp instance
             */
            public static create(properties?: google.protobuf.ITimestamp): google.protobuf.Timestamp;

            /**
             * Encodes the specified Timestamp message. Does not implicitly {@link google.protobuf.Timestamp.verify|verify} messages.
             * @param message Timestamp message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encode(message: google.protobuf.ITimestamp, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Encodes the specified Timestamp message, length delimited. Does not implicitly {@link google.protobuf.Timestamp.verify|verify} messages.
             * @param message Timestamp message or plain object to encode
             * @param [writer] Writer to encode to
             * @returns Writer
             */
            public static encodeDelimited(message: google.protobuf.ITimestamp, writer?: $protobuf.Writer): $protobuf.Writer;

            /**
             * Decodes a Timestamp message from the specified reader or buffer.
             * @param reader Reader or buffer to decode from
             * @param [length] Message length if known beforehand
             * @returns Timestamp
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decode(reader: ($protobuf.Reader|Uint8Array), length?: number): google.protobuf.Timestamp;

            /**
             * Decodes a Timestamp message from the specified reader or buffer, length delimited.
             * @param reader Reader or buffer to decode from
             * @returns Timestamp
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            public static decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): google.protobuf.Timestamp;

            /**
             * Verifies a Timestamp message.
             * @param message Plain object to verify
             * @returns `null` if valid, otherwise the reason why it is not
             */
            public static verify(message: { [k: string]: any }): (string|null);

            /**
             * Creates a Timestamp message from a plain object. Also converts values to their respective internal types.
             * @param object Plain object
             * @returns Timestamp
             */
            public static fromObject(object: { [k: string]: any }): google.protobuf.Timestamp;

            /**
             * Creates a plain object from a Timestamp message. Also converts values to other types if specified.
             * @param message Timestamp
             * @param [options] Conversion options
             * @returns Plain object
             */
            public static toObject(message: google.protobuf.Timestamp, options?: $protobuf.IConversionOptions): { [k: string]: any };

            /**
             * Converts this Timestamp to JSON.
             * @returns JSON object
             */
            public toJSON(): { [k: string]: any };

            /**
             * Gets the default type url for Timestamp
             * @param [typeUrlPrefix] your custom typeUrlPrefix(default "type.googleapis.com")
             * @returns The default type url
             */
            public static getTypeUrl(typeUrlPrefix?: string): string;
        }
    }
}
