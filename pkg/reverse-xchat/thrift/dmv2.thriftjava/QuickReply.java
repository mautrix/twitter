package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.NoWhenBranchMatchedException;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.DefaultConstructorMarker;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000(\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u0007\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\b6\u0018\u0000 \t2\u00020\u0001:\u0005\n\u000b\f\r\tB\t\b\u0004¢\u0006\u0004\b\u0002\u0010\u0003J\u0017\u0010\u0007\u001a\u00020\u00062\u0006\u0010\u0005\u001a\u00020\u0004H\u0016¢\u0006\u0004\b\u0007\u0010\b\u0082\u0001\u0003\u000e\u000f\u0010¨\u0006\u0011"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/QuickReply;", "Lcom/bendb/thrifty/a;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "Companion", "Request", "Response", "Unknown", "QuickReplyAdapter", "Lcom/x/dmv2/thriftjava/QuickReply$Request;", "Lcom/x/dmv2/thriftjava/QuickReply$Response;", "Lcom/x/dmv2/thriftjava/QuickReply$Unknown;", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public abstract class QuickReply implements InterfaceC11261a {

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new QuickReplyAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/QuickReply$QuickReplyAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/QuickReply;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/QuickReply;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/QuickReply;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class QuickReplyAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public QuickReply m85952read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            QuickReply response;
            Intrinsics.m65272h(protocol, "protocol");
            QuickReply quickReply = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    break;
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        quickReply = Unknown.INSTANCE;
                        C11272a.m14141a(protocol, b);
                    } else if (b == 12) {
                        response = new Response((QuickReplyResponse) QuickReplyResponse.ADAPTER.read(protocol));
                        quickReply = response;
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 12) {
                    response = new Request((QuickReplyRequest) QuickReplyRequest.ADAPTER.read(protocol));
                    quickReply = response;
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
            if (quickReply != null) {
                return quickReply;
            }
            throw new IllegalStateException("unreadable");
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a QuickReply struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("QuickReply");
            if (struct instanceof Request) {
                protocol.mo14136v3("request", 1, (byte) 12);
                QuickReplyRequest.ADAPTER.write(protocol, ((Request) struct).m76803getValue());
            } else if (struct instanceof Response) {
                protocol.mo14136v3("response", 2, (byte) 12);
                QuickReplyResponse.ADAPTER.write(protocol, ((Response) struct).m76804getValue());
            } else if (!(struct instanceof Unknown)) {
                throw new NoWhenBranchMatchedException();
            }
            protocol.mo14134i0();
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/QuickReply$Request;", "Lcom/x/dmv2/thriftjava/QuickReply;", "value", "Lcom/x/dmv2/thriftjava/QuickReplyRequest;", "<init>", "(Lcom/x/dmv2/thriftjava/QuickReplyRequest;)V", "getValue", "()Lcom/x/dmv2/thriftjava/QuickReplyRequest;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Request extends QuickReply {

        @InterfaceC88464a
        private final QuickReplyRequest value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public Request(@InterfaceC88464a QuickReplyRequest value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ Request copy$default(Request request, QuickReplyRequest quickReplyRequest, int i, Object obj) {
            if ((i & 1) != 0) {
                quickReplyRequest = request.value;
            }
            return request.copy(quickReplyRequest);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final QuickReplyRequest getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final Request copy(@InterfaceC88464a QuickReplyRequest value) {
            Intrinsics.m65272h(value, "value");
            return new Request(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof Request) && Intrinsics.m65267c(this.value, ((Request) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final QuickReplyRequest m76803getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "QuickReply(request=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/QuickReply$Response;", "Lcom/x/dmv2/thriftjava/QuickReply;", "value", "Lcom/x/dmv2/thriftjava/QuickReplyResponse;", "<init>", "(Lcom/x/dmv2/thriftjava/QuickReplyResponse;)V", "getValue", "()Lcom/x/dmv2/thriftjava/QuickReplyResponse;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Response extends QuickReply {

        @InterfaceC88464a
        private final QuickReplyResponse value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public Response(@InterfaceC88464a QuickReplyResponse value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ Response copy$default(Response response, QuickReplyResponse quickReplyResponse, int i, Object obj) {
            if ((i & 1) != 0) {
                quickReplyResponse = response.value;
            }
            return response.copy(quickReplyResponse);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final QuickReplyResponse getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final Response copy(@InterfaceC88464a QuickReplyResponse value) {
            Intrinsics.m65272h(value, "value");
            return new Response(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof Response) && Intrinsics.m65267c(this.value, ((Response) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final QuickReplyResponse m76804getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "QuickReply(response=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000$\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\n\u0002\u0010\u000e\n\u0000\bÆ\n\u0018\u00002\u00020\u0001B\t\b\u0002¢\u0006\u0004\b\u0002\u0010\u0003J\u0013\u0010\u0004\u001a\u00020\u00052\b\u0010\u0006\u001a\u0004\u0018\u00010\u0007HÖ\u0003J\t\u0010\b\u001a\u00020\tHÖ\u0001J\t\u0010\n\u001a\u00020\u000bHÖ\u0001¨\u0006\f"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/QuickReply$Unknown;", "Lcom/x/dmv2/thriftjava/QuickReply;", "<init>", "()V", "equals", "", "other", "", "hashCode", "", "toString", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Unknown extends QuickReply {

        @InterfaceC88464a
        public static final Unknown INSTANCE = new Unknown();

        private Unknown() {
            super(null);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            return this == other || (other instanceof Unknown);
        }

        public int hashCode() {
            return -1842477082;
        }

        @InterfaceC88464a
        public String toString() {
            return "Unknown";
        }
    }

    public /* synthetic */ QuickReply(DefaultConstructorMarker defaultConstructorMarker) {
        this();
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }

    private QuickReply() {
    }
}