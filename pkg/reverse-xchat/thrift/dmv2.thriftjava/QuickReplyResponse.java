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

@Metadata(m64929d1 = {"\u0000$\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u0006\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\b6\u0018\u0000 \t2\u00020\u0001:\u0004\n\u000b\f\tB\t\b\u0004¢\u0006\u0004\b\u0002\u0010\u0003J\u0017\u0010\u0007\u001a\u00020\u00062\u0006\u0010\u0005\u001a\u00020\u0004H\u0016¢\u0006\u0004\b\u0007\u0010\b\u0082\u0001\u0002\r\u000e¨\u0006\u000f"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/QuickReplyResponse;", "Lcom/bendb/thrifty/a;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "Companion", "Options", "Unknown", "QuickReplyResponseAdapter", "Lcom/x/dmv2/thriftjava/QuickReplyResponse$Options;", "Lcom/x/dmv2/thriftjava/QuickReplyResponse$Unknown;", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public abstract class QuickReplyResponse implements InterfaceC11261a {

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new QuickReplyResponseAdapter();

    @Metadata(m64929d1 = {"\u0000,\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0005\n\u0002\u0010\u000e\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\b\u0086\b\u0018\u00002\u00020\u0001B\u000f\u0012\u0006\u0010\u0002\u001a\u00020\u0003¢\u0006\u0004\b\u0004\u0010\u0005J\b\u0010\b\u001a\u00020\tH\u0016J\t\u0010\n\u001a\u00020\u0003HÆ\u0003J\u0013\u0010\u000b\u001a\u00020\u00002\b\b\u0002\u0010\u0002\u001a\u00020\u0003HÆ\u0001J\u0013\u0010\f\u001a\u00020\r2\b\u0010\u000e\u001a\u0004\u0018\u00010\u000fHÖ\u0003J\t\u0010\u0010\u001a\u00020\u0011HÖ\u0001R\u0011\u0010\u0002\u001a\u00020\u0003¢\u0006\b\n\u0000\u001a\u0004\b\u0006\u0010\u0007¨\u0006\u0012"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/QuickReplyResponse$Options;", "Lcom/x/dmv2/thriftjava/QuickReplyResponse;", "value", "Lcom/x/dmv2/thriftjava/QuickReplyOptionsResponse;", "<init>", "(Lcom/x/dmv2/thriftjava/QuickReplyOptionsResponse;)V", "getValue", "()Lcom/x/dmv2/thriftjava/QuickReplyOptionsResponse;", "toString", "", "component1", "copy", "equals", "", "other", "", "hashCode", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Options extends QuickReplyResponse {

        @InterfaceC88464a
        private final QuickReplyOptionsResponse value;

        /* JADX WARN: 'super' call moved to the top of the method (can break code semantics) */
        public Options(@InterfaceC88464a QuickReplyOptionsResponse value) {
            super(null);
            Intrinsics.m65272h(value, "value");
            this.value = value;
        }

        public static /* synthetic */ Options copy$default(Options options, QuickReplyOptionsResponse quickReplyOptionsResponse, int i, Object obj) {
            if ((i & 1) != 0) {
                quickReplyOptionsResponse = options.value;
            }
            return options.copy(quickReplyOptionsResponse);
        }

        @InterfaceC88464a
        /* renamed from: component1, reason: from getter */
        public final QuickReplyOptionsResponse getValue() {
            return this.value;
        }

        @InterfaceC88464a
        public final Options copy(@InterfaceC88464a QuickReplyOptionsResponse value) {
            Intrinsics.m65272h(value, "value");
            return new Options(value);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            if (this == other) {
                return true;
            }
            return (other instanceof Options) && Intrinsics.m65267c(this.value, ((Options) other).value);
        }

        @InterfaceC88464a
        /* renamed from: getValue */
        public final QuickReplyOptionsResponse m76806getValue() {
            return this.value;
        }

        public int hashCode() {
            return this.value.hashCode();
        }

        @InterfaceC88464a
        public String toString() {
            return "QuickReplyResponse(options=" + this.value + Separators.RPAREN;
        }
    }

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/QuickReplyResponse$QuickReplyResponseAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/QuickReplyResponse;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/QuickReplyResponse;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/QuickReplyResponse;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class QuickReplyResponseAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public QuickReplyResponse m85954read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            QuickReplyResponse options = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    break;
                }
                if (c11265cMo14127V2.f38393b != 1) {
                    options = Unknown.INSTANCE;
                    C11272a.m14141a(protocol, b);
                } else if (b == 12) {
                    options = new Options((QuickReplyOptionsResponse) QuickReplyOptionsResponse.ADAPTER.read(protocol));
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
            if (options != null) {
                return options;
            }
            throw new IllegalStateException("unreadable");
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a QuickReplyResponse struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("QuickReplyResponse");
            if (struct instanceof Options) {
                protocol.mo14136v3("options", 1, (byte) 12);
                QuickReplyOptionsResponse.ADAPTER.write(protocol, ((Options) struct).m76806getValue());
            } else if (!(struct instanceof Unknown)) {
                throw new NoWhenBranchMatchedException();
            }
            protocol.mo14134i0();
        }
    }

    @Metadata(m64929d1 = {"\u0000$\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0010\u000b\n\u0000\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\b\n\u0000\n\u0002\u0010\u000e\n\u0000\bÆ\n\u0018\u00002\u00020\u0001B\t\b\u0002¢\u0006\u0004\b\u0002\u0010\u0003J\u0013\u0010\u0004\u001a\u00020\u00052\b\u0010\u0006\u001a\u0004\u0018\u00010\u0007HÖ\u0003J\t\u0010\b\u001a\u00020\tHÖ\u0001J\t\u0010\n\u001a\u00020\u000bHÖ\u0001¨\u0006\f"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/QuickReplyResponse$Unknown;", "Lcom/x/dmv2/thriftjava/QuickReplyResponse;", "<init>", "()V", "equals", "", "other", "", "hashCode", "", "toString", "", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final /* data */ class Unknown extends QuickReplyResponse {

        @InterfaceC88464a
        public static final Unknown INSTANCE = new Unknown();

        private Unknown() {
            super(null);
        }

        public boolean equals(@InterfaceC88465b Object other) {
            return this == other || (other instanceof Unknown);
        }

        public int hashCode() {
            return 1757693255;
        }

        @InterfaceC88464a
        public String toString() {
            return "Unknown";
        }
    }

    public /* synthetic */ QuickReplyResponse(DefaultConstructorMarker defaultConstructorMarker) {
        this();
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }

    private QuickReplyResponse() {
    }
}